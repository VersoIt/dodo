package kitchen

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type TicketStatus int

const (
	TicketQueued  TicketStatus = 0
	TicketCooking TicketStatus = 1
	TicketReady   TicketStatus = 2
)

func (s TicketStatus) String() string {
	switch s {
	case TicketQueued:
		return "queued"
	case TicketCooking:
		return "cooking"
	case TicketReady:
		return "ready"
	default:
		return "unknown"
	}
}

type KitchenTicket struct {
	id               string
	orderID          string
	items            []KitchenItem
	status           TicketStatus
	createdAt        time.Time
	startCookingTime time.Time
	readyTime        time.Time
}

type KitchenItem struct {
	ProductID   string
	Name        string
	Ingredients []string
	Quantity    int
	Comment     string
}

func NewTicket(orderID string, items []KitchenItem) *KitchenTicket {
	id, _ := uuid.NewV7()
	return &KitchenTicket{
		id:        id.String(),
		orderID:   orderID,
		items:     items,
		status:    TicketQueued,
		createdAt: time.Now(),
	}
}

func (t *KitchenTicket) StartCooking() error {
	if t.status != TicketQueued {
		return errors.New("ticket is not in queue")
	}
	t.status = TicketCooking
	t.startCookingTime = time.Now()
	return nil
}

func (t *KitchenTicket) MarkReady() error {
	if t.status != TicketCooking {
		return errors.New("ticket must be cooking before ready")
	}
	t.status = TicketReady
	t.readyTime = time.Now()
	return nil
}

func (t *KitchenTicket) GetCookingDuration() time.Duration {
	if t.readyTime.IsZero() || t.startCookingTime.IsZero() {
		return 0
	}
	return t.readyTime.Sub(t.startCookingTime)
}

func (t *KitchenTicket) ID() string           { return t.id }
func (t *KitchenTicket) OrderID() string      { return t.orderID }
func (t *KitchenTicket) Status() TicketStatus { return t.status }
func (t *KitchenTicket) Items() []KitchenItem { return t.items }
func (t *KitchenTicket) CreatedAt() time.Time { return t.createdAt }
func (t *KitchenTicket) StartTime() time.Time { return t.startCookingTime }
func (t *KitchenTicket) ReadyTime() time.Time { return t.readyTime }

type TicketRepository interface {
	Save(ctx context.Context, t *KitchenTicket) error
	FindPending(ctx context.Context) ([]*KitchenTicket, error)
	FindByID(ctx context.Context, id string) (*KitchenTicket, error)
}