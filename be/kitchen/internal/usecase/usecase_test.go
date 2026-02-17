package usecase

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/versoit/diploma/be/kitchen/internal/domain"
)

type MockTicketRepo struct {
	store map[string]*domain.KitchenTicket
}

func (m *MockTicketRepo) Save(ctx context.Context, t *domain.KitchenTicket) error {
	m.store[t.ID()] = t
	return nil
}

func (m *MockTicketRepo) FindByID(ctx context.Context, id string) (*domain.KitchenTicket, error) {
	if t, ok := m.store[id]; ok {
		return t, nil
	}
	return nil, errors.New("not found")
}

func (m *MockTicketRepo) FindPending(ctx context.Context) ([]*domain.KitchenTicket, error) {
	return nil, nil
}

func (m *MockTicketRepo) FindAll(ctx context.Context) ([]*domain.KitchenTicket, error) {
	return nil, nil
}

type dummyTM struct{}

func (d *dummyTM) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func (d *dummyTM) DoWithSettings(ctx context.Context, _ trm.Settings, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestKitchenUseCase_AcceptOrder(t *testing.T) {
	repo := &MockTicketRepo{store: make(map[string]*domain.KitchenTicket)}
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	uc := NewKitchenUseCase(repo, &dummyTM{}, log)

	items := []domain.KitchenItem{{Name: "Pizza", Quantity: 1}}
	ticket, err := uc.AcceptOrder(context.Background(), "order-123", "ORD-123", items)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ticket.OrderID() != "order-123" {
		t.Errorf("expected orderID order-123, got %s", ticket.OrderID())
	}
}

func TestKitchenUseCase_CookingFlow(t *testing.T) {
	repo := &MockTicketRepo{store: make(map[string]*domain.KitchenTicket)}
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	uc := NewKitchenUseCase(repo, &dummyTM{}, log)

	ticket, err := uc.AcceptOrder(context.Background(), "ord-1", "N-1", []domain.KitchenItem{{Name: "P", Quantity: 1}})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	id := ticket.ID()

	_, err = uc.StartCooking(context.Background(), id)
	if err != nil {
		t.Fatalf("start cooking failed: %v", err)
	}

	saved, _ := repo.FindByID(context.Background(), id)
	if saved.Status() != domain.TicketCooking {
		t.Errorf("expected cooking status, got %v", saved.Status())
	}
}
