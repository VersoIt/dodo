package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleClient  Role = "client"
	RoleCourier Role = "courier"
	RoleSupport Role = "support"
	RoleSystem  Role = "system"
)

type Message struct {
	ID         int64     `json:"message_id"`
	OrderID    uuid.UUID `json:"order_id"`
	SenderID   uuid.UUID `json:"sender_id"` // User ID from Auth
	SenderName string    `json:"sender_name"`
	Role       Role      `json:"role"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"created_at"`
	IsRead     bool      `json:"is_read"`
}

type MessageRepository interface {
	Save(ctx context.Context, msg *Message) error
	GetHistory(ctx context.Context, orderID uuid.UUID, limit int) ([]Message, error)
	GetAfterID(ctx context.Context, orderID uuid.UUID, afterID int64) ([]Message, error)
	MarkAsRead(ctx context.Context, messageID int64) error
}

// WS Events

type WSEventType string

const (
	EventSendMessage WSEventType = "send_message"
	EventMessageAck  WSEventType = "message_ack"
	EventNewMessage  WSEventType = "new_message"
	EventRead        WSEventType = "read"
	EventError       WSEventType = "error"
)

type WSMessage struct {
	Action     WSEventType `json:"action,omitempty"` // From Client
	Event      WSEventType `json:"event,omitempty"`  // To Client
	RequestID  string      `json:"request_id,omitempty"`
	MessageID  int64       `json:"message_id,omitempty"`
	Text       string      `json:"text,omitempty"`
	Role       Role        `json:"role,omitempty"`
	SenderName string      `json:"sender_name,omitempty"`
	OrderID    uuid.UUID   `json:"order_id,omitempty"`
	CreatedAt  time.Time   `json:"created_at,omitempty"`
	Payload    interface{} `json:"payload,omitempty"`
}
