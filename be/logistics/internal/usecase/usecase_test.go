package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/versoit/diploma/be/logistics/internal/domain"
)

type MockDeliveryRepo struct {
	store map[string]*domain.Delivery
}

func (m *MockDeliveryRepo) Save(ctx context.Context, d *domain.Delivery) error {
	m.store[d.OrderID()] = d
	return nil
}
func (m *MockDeliveryRepo) FindByOrderID(ctx context.Context, id string) (*domain.Delivery, error) {
	if d, ok := m.store[id]; ok {
		return d, nil
	}
	return nil, fmt.Errorf("not found")
}
func (m *MockDeliveryRepo) FindAll(ctx context.Context) ([]*domain.Delivery, error) {
	return nil, nil
}

type MockCourierRepo struct {
	store map[string]*domain.Courier
}

func (m *MockCourierRepo) Save(ctx context.Context, c *domain.Courier) error {
	m.store[c.ID()] = c
	return nil
}
func (m *MockCourierRepo) FindByID(ctx context.Context, id string) (*domain.Courier, error) {
	if c, ok := m.store[id]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("not found")
}
func (m *MockCourierRepo) FindAvailable(ctx context.Context) ([]*domain.Courier, error) {
	return nil, nil
}
func (m *MockCourierRepo) UpdateLocation(ctx context.Context, id string, lat, lng float64) error {
	return nil
}

type dummyTM struct{}

func (d *dummyTM) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func (d *dummyTM) DoWithSettings(ctx context.Context, _ trm.Settings, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestLogisticsUseCase_AssignCourier(t *testing.T) {
	dRepo := &MockDeliveryRepo{store: make(map[string]*domain.Delivery)}
	cRepo := &MockCourierRepo{store: make(map[string]*domain.Courier)}
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	uc := NewLogisticsUseCase(dRepo, cRepo, &dummyTM{}, log)

	courier := domain.NewCourier("Vasya", "123")
	courier.GoOnline()
	if err := cRepo.Save(context.Background(), courier); err != nil {
		t.Fatalf("failed to save courier: %v", err)
	}

	// Create a dummy delivery first
	delivery := domain.NewDelivery("order-1", "ORD-1", "City", "Street", "1", "10", nil)
	if err := dRepo.Save(context.Background(), delivery); err != nil {
		t.Fatalf("failed to save delivery: %v", err)
	}

	err := uc.AssignCourierToDelivery(context.Background(), "order-1", courier.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	d, err := dRepo.FindByOrderID(context.Background(), "order-1")
	if err != nil {
		t.Fatalf("failed to find delivery: %v", err)
	}
	if d.CourierID() != courier.ID() {
		t.Error("courier not assigned")
	}
}
