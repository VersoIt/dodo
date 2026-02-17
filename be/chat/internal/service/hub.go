package service

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/versoit/diploma/be/chat/internal/domain"
)

type Client struct {
	Conn    *websocket.Conn
	Hub     *Hub
	OrderID uuid.UUID
	Send    chan []byte
}

type Hub struct {
	// Registered clients.
	Clients map[uuid.UUID]map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan domain.WSMessage

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	natsConn *nats.Conn
	log      *slog.Logger
	mu       sync.RWMutex
}

func NewHub(nc *nats.Conn, log *slog.Logger) *Hub {
	return &Hub{
		Broadcast:  make(chan domain.WSMessage),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[uuid.UUID]map[*Client]bool),
		natsConn:   nc,
		log:        log,
	}
}

func (h *Hub) Run() {
	// Subscribe to NATS for cross-instance communication
	// Subject: chat.order.*
	if h.natsConn != nil {
		_, err := h.natsConn.Subscribe("chat.order.*", func(m *nats.Msg) {
			var msg domain.WSMessage
			if err := json.Unmarshal(m.Data, &msg); err != nil {
				h.log.Error("Failed to unmarshal NATS message", slog.Any("error", err))
				return
			}
			h.broadcastLocal(msg)
		})
		if err != nil {
			h.log.Error("Failed to subscribe to NATS", slog.Any("error", err))
		}
	}

	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.Clients[client.OrderID] == nil {
				h.Clients[client.OrderID] = make(map[*Client]bool)
			}
			h.Clients[client.OrderID][client] = true
			h.mu.Unlock()
			h.log.Debug("Client registered", slog.String("order_id", client.OrderID.String()))

		case client := <-h.Unregister:
			h.mu.Lock()
			if clients, ok := h.Clients[client.OrderID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.Clients, client.OrderID)
					}
				}
			}
			h.mu.Unlock()
			h.log.Debug("Client unregistered", slog.String("order_id", client.OrderID.String()))

		case message := <-h.Broadcast:
			// Publish to NATS to distribute to other instances
			if h.natsConn != nil {
				data, err := json.Marshal(message)
				if err == nil {
					subject := fmt.Sprintf("chat.order.%s", message.OrderID.String())
					if err := h.natsConn.Publish(subject, data); err != nil {
						h.log.Error("Failed to publish to NATS", slog.Any("error", err))
					}
				} else {
					h.log.Error("Failed to marshal WS message for NATS", slog.Any("error", err))
				}
			}
			// Also broadcast locally immediately
			h.broadcastLocal(message)
		}
	}
}

func (h *Hub) broadcastLocal(message domain.WSMessage) {
	h.mu.RLock()
	clients, ok := h.Clients[message.OrderID]
	if !ok {
		h.mu.RUnlock()
		return
	}

	data, err := json.Marshal(message)
	if err != nil {
		h.mu.RUnlock()
		h.log.Error("Failed to marshal WS message", slog.Any("error", err))
		return
	}

	// Iterate over clients under read lock
	for client := range clients {
		select {
		case client.Send <- data:
		default:
		}
	}
	h.mu.RUnlock()
}
