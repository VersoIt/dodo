package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	pb "github.com/versoit/diploma/be/chat/api/proto/pb"
)

type ChatHandler struct {
	client pb.ChatServiceClient
	log    *slog.Logger
	// We need the JWT secret here to validate the token from Query param
	// Usually this is in config, but for now I'll assume it's passed or available.
	// Ideally Gateway config has it.
}

func NewChatHandler(client pb.ChatServiceClient, log *slog.Logger) *ChatHandler {
	return &ChatHandler{
		client: client,
		log:    log,
	}
}

// WSUpgrade handles the WebSocket upgrade
func (h *ChatHandler) WSUpgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

// HandleWS manages the WebSocket connection and gRPC stream
func (h *ChatHandler) HandleWS(c *websocket.Conn) {
	// 1. Authenticate (Token in Query or Protocol)
	tokenStr := c.Query("token")
	if tokenStr == "" {
		tokenStr = c.Cookies("token") // or header
	}
	// Simplified validation: In a real gateway, we'd use a shared secret or call Auth service.
	// Assuming "super-secret-key" for now as per other services, BUT Gateway should probably get user info from Auth middleware if possible.
	// Since WS upgrade is tricky with middleware, we parse manually.
	
	// FIX: Hardcoded secret for demo purposes, should be from Config.
	token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte("super-secret-key"), nil
	})

	var userID, role string
	if token != nil && token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			userID, _ = claims["user_id"].(string)
			role, _ = claims["role"].(string)
		}
	}

	if userID == "" {
		_ = c.WriteJSON(fiber.Map{"error": "Unauthorized"})
		_ = c.Close()
		return
	}

	orderID := c.Query("order_id")
	if orderID == "" {
		_ = c.WriteJSON(fiber.Map{"error": "order_id required"})
		_ = c.Close()
		return
	}

	// 2. Connect to gRPC
	stream, err := h.client.Connect(context.Background())
	if err != nil {
		h.log.Error("Failed to connect to Chat Service", slog.Any("error", err))
		_ = c.Close()
		return
	}

	// 3. Send Handshake
	err = stream.Send(&pb.ClientMessage{
		Event: &pb.ClientMessage_Connect{
			Connect: &pb.ConnectEvent{
				UserId:  userID,
				Role:    role,
				OrderId: orderID,
			},
		},
	})
	if err != nil {
		h.log.Error("Failed to send handshake", slog.Any("error", err))
		_ = c.Close()
		return
	}

	// 4. Bidirectional Proxy
	var wg sync.WaitGroup
	wg.Add(2)

	done := make(chan struct{})

	// Goroutine: WS -> gRPC
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
				_, msg, err := c.ReadMessage()
				if err != nil {
					h.log.Debug("WS read closed", slog.Any("error", err))
					_ = stream.CloseSend()
					close(done)
					return
				}

				var wsReq struct {
					Action     string `json:"action"`
					Text       string `json:"text"`
					RequestID  string `json:"request_id"`
					MessageID  int64  `json:"message_id"`
				}

				if err := json.Unmarshal(msg, &wsReq); err != nil {
					continue
				}

				var pbMsg *pb.ClientMessage
				switch wsReq.Action {
				case "send_message":
					pbMsg = &pb.ClientMessage{
						Event: &pb.ClientMessage_SendMessage{
							SendMessage: &pb.SendMessageEvent{
								RequestId: wsReq.RequestID,
								Text:      wsReq.Text,
							},
						},
					}
				case "read_message":
					pbMsg = &pb.ClientMessage{
						Event: &pb.ClientMessage_ReadMessage{
							ReadMessage: &pb.ReadMessageEvent{
								MessageId: wsReq.MessageID,
							},
						},
					}
				}

				if pbMsg != nil {
					if err := stream.Send(pbMsg); err != nil {
						h.log.Error("gRPC send error", slog.Any("error", err))
						return
					}
				}
			}
		}
	}()

	// Goroutine: gRPC -> WS
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
				in, err := stream.Recv()
				if err == io.EOF {
					return
				}
				if err != nil {
					h.log.Error("gRPC recv error", slog.Any("error", err))
					_ = c.Close()
					close(done)
					return
				}

				var wsResp map[string]interface{}

				switch event := in.Event.(type) {
				case *pb.ServerMessage_NewMessage:
					wsResp = map[string]interface{}{
						"event":      "new_message",
						"message_id": event.NewMessage.MessageId,
						"text":       event.NewMessage.Text,
						"role":       event.NewMessage.Role,
						"order_id":   event.NewMessage.OrderId,
						"created_at": event.NewMessage.CreatedAt.AsTime(),
					}
				case *pb.ServerMessage_Ack:
					wsResp = map[string]interface{}{
						"event":      "message_ack",
						"request_id": event.Ack.RequestId,
						"message_id": event.Ack.MessageId,
					}
				case *pb.ServerMessage_Error:
					wsResp = map[string]interface{}{
						"event":   "error",
						"payload": event.Error.Message,
					}
				}

				if wsResp != nil {
					if err := c.WriteJSON(wsResp); err != nil {
						h.log.Error("WS write error", slog.Any("error", err))
						return
					}
				}
			}
		}
	}()

	wg.Wait()
}
