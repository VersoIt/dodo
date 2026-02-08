package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/versoit/diploma/services/logistics"
)

type MockDeliveryRepo struct {
	store map[string]*logistics.Delivery
}

func (m *MockDeliveryRepo) Save(ctx context.Context, d *logistics.Delivery) error {
	m.store[d.OrderID()] = d
	return nil
}
func (m *MockDeliveryRepo) FindByOrderID(ctx context.Context, id string) (*logistics.Delivery, error) {
	if d, ok := m.store[id]; ok {
		return d, nil
	}
	return nil, fmt.Errorf("not found")
}

type MockCourierRepo struct {
	store map[string]*logistics.Courier
}

func (m *MockCourierRepo) Save(ctx context.Context, c *logistics.Courier) error {
	m.store[c.ID()] = c
	return nil
}
func (m *MockCourierRepo) FindByID(ctx context.Context, id string) (*logistics.Courier, error) {
	if c, ok := m.store[id]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("not found")
}
func (m *MockCourierRepo) FindAvailable(ctx context.Context) ([]*logistics.Courier, error) {
	return nil, nil
}
func (m *MockCourierRepo) UpdateLocation(ctx context.Context, id string, lat, lng float64) error {
	return nil
}

func TestLogisticsUseCase_AssignCourier(t *testing.T) {
	dRepo := &MockDeliveryRepo{store: make(map[string]*logistics.Delivery)}
	cRepo := &MockCourierRepo{store: make(map[string]*logistics.Courier)}
	uc := NewLogisticsUseCase(dRepo, cRepo)

	courier := logistics.NewCourier("Vasya", "123")
	courier.GoOnline()
	if err := cRepo.Save(context.Background(), courier); err != nil {
		t.Fatalf("failed to save courier: %v", err)
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
