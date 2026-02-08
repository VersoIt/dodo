package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/versoit/diploma/services/kitchen"
)

type MockTicketRepo struct {
	store map[string]*kitchen.KitchenTicket
}

func (m *MockTicketRepo) Save(ctx context.Context, t *kitchen.KitchenTicket) error {
	m.store[t.ID()] = t
	return nil
}

func (m *MockTicketRepo) FindByID(ctx context.Context, id string) (*kitchen.KitchenTicket, error) {
	if t, ok := m.store[id]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("ticket not found")
}

func (m *MockTicketRepo) FindPending(ctx context.Context) ([]*kitchen.KitchenTicket, error) {
	return nil, nil
}

func TestKitchenUseCase_AcceptOrder(t *testing.T) {
	repo := &MockTicketRepo{store: make(map[string]*kitchen.KitchenTicket)}
	uc := NewKitchenUseCase(repo)

	items := []kitchen.KitchenItem{{Name: "Pizza", Quantity: 1}}
	ticket, err := uc.AcceptOrder(context.Background(), "order-123", items)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ticket.OrderID() != "order-123" {
		t.Errorf("expected orderID order-123, got %s", ticket.OrderID())
	}
}

func TestKitchenUseCase_CookingFlow(t *testing.T) {
	repo := &MockTicketRepo{store: make(map[string]*kitchen.KitchenTicket)}
	uc := NewKitchenUseCase(repo)

	ticket, err := uc.AcceptOrder(context.Background(), "ord-1", []kitchen.KitchenItem{{Name: "P"}})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	id := ticket.ID()

	err = uc.StartCooking(context.Background(), id)
	if err != nil {
		t.Fatalf("start cooking failed: %v", err)
	}

	saved, err := repo.FindByID(context.Background(), id)
	if err != nil {
		t.Fatalf("failed to find ticket: %v", err)
	}
	if saved.Status() != kitchen.TicketCooking {
		t.Errorf("expected cooking status, got %v", saved.Status())
	}
}
