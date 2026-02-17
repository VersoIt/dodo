package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"strconv"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/versoit/diploma/be/chat/internal/config"
	"github.com/versoit/diploma/be/chat/internal/domain"
	"github.com/versoit/diploma/be/chat/internal/service"
)

type ChatHandler struct {
	svc service.MessageService
	hub *service.Hub
	cfg *config.Config
	log *slog.Logger
}

func NewChatHandler(svc service.MessageService, hub *service.Hub, cfg *config.Config, log *slog.Logger) *ChatHandler {
	return &ChatHandler{
		svc: svc,
		hub: hub,
		cfg: cfg,
		log: log,
	}
}

func (h *ChatHandler) GetHistory(c *fiber.Ctx) error {
	orderIDStr := c.Query("order_id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid order_id (UUID required)"})
	}

	limitStr := c.Query("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	messages, err := h.svc.GetHistory(c.Context(), orderID, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(messages)
}

func (h *ChatHandler) GetSync(c *fiber.Ctx) error {
	orderIDStr := c.Query("order_id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid order_id (UUID required)"})
	}

	afterIDStr := c.Query("after_id")
	afterID, _ := strconv.ParseInt(afterIDStr, 10, 64)

	messages, err := h.svc.GetUpdates(c.Context(), orderID, afterID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(messages)
}

func (h *ChatHandler) WSUpgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			tokenStr = c.Get("Sec-WebSocket-Protocol")
		}

		if tokenStr == "" {
			return fiber.ErrUnauthorized
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return fiber.ErrUnauthorized
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return fiber.ErrUnauthorized
		}

		c.Locals("user_id", claims["user_id"])
		c.Locals("role", claims["role"])

		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func (h *ChatHandler) WebSocketHandler(c *websocket.Conn) {
	userIDStr, _ := c.Locals("user_id").(string)
	roleStr, _ := c.Locals("role").(string)
	orderIDStr := c.Query("order_id")
	orderID, err := uuid.Parse(orderIDStr)

	if err != nil {
		h.log.Warn("WS connection with invalid order_id", slog.String("user_id", userIDStr), slog.String("order_id", orderIDStr))
		_ = c.WriteJSON(fiber.Map{"error": "order_id query param is required and must be UUID"})
		_ = c.Close()
		return
	}

	userID, _ := uuid.Parse(userIDStr)

	client := &service.Client{
		Conn:    c,
		Hub:     h.hub,
		OrderID: orderID,
		Send:    make(chan []byte, 256),
	}
	h.hub.Register <- client

	go func() {
		for msg := range client.Send {
			if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
				h.log.Error("WS write error", slog.Any("error", err))
				break
			}
		}
	}()

	defer func() {
		h.hub.Unregister <- client
		_ = c.Close()
	}()

	for {
		_, payload, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.log.Error("WS read error", slog.Any("error", err))
			}
			break
		}

		var wsReq domain.WSMessage
		if err := json.Unmarshal(payload, &wsReq); err != nil {
			h.log.Warn("Invalid WS payload", slog.Any("error", err))
			continue
		}

		switch wsReq.Action {
		case domain.EventSendMessage:
			msg := &domain.Message{
				OrderID:  orderID,
				SenderID: userID,
				Role:     domain.Role(roleStr),
				Text:     wsReq.Text,
			}
			
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err := h.svc.Save(ctx, msg)
			cancel()

			if err != nil {
				h.log.Error("Failed to save message", slog.Any("error", err))
				_ = c.WriteJSON(domain.WSMessage{
					Event:   domain.EventError,
					Payload: "Failed to save message",
				})
				continue
			}

			_ = c.WriteJSON(domain.WSMessage{
				Event:     domain.EventMessageAck,
				RequestID: wsReq.RequestID,
				MessageID: msg.ID,
			})

			h.svc.Broadcast(domain.WSMessage{
				Event:     domain.EventNewMessage,
				MessageID: msg.ID,
				Text:      msg.Text,
				Role:      msg.Role,
				OrderID:   orderID,
				CreatedAt: msg.CreatedAt,
			})

		case domain.EventRead:
			if wsReq.MessageID != 0 {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				_ = h.svc.MarkAsRead(ctx, wsReq.MessageID)
				cancel()
			}
		}
	}
}
