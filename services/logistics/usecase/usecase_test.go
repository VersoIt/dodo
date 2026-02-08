package usecase

import (
	"context"
	"testing"
	"fmt"

	"github.com/versoit/diploma/services/logistics"
)

// Mock Repos
type MockDeliveryRepo struct {
	store map[string]*logistics.Delivery
}
func (m *MockDeliveryRepo) Save(d *logistics.Delivery) error {
	// Key by OrderID for simplicity based on FindByOrderID usage
	m.store[d.OrderID()] = d 
	return nil
}
func (m *MockDeliveryRepo) FindByOrderID(id string) (*logistics.Delivery, error) {
	if d, ok := m.store[id]; ok {
		return d, nil
	}
	return nil, fmt.Errorf("not found")
}

type MockCourierRepo struct {
	store map[string]*logistics.Courier
}
func (m *MockCourierRepo) Save(c *logistics.Courier) error {
	m.store[c.ID()] = c
	return nil
}
func (m *MockCourierRepo) FindByID(id string) (*logistics.Courier, error) {
	if c, ok := m.store[id]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("not found")
}
func (m *MockCourierRepo) FindAvailable() ([]*logistics.Courier, error) { return nil, nil }
func (m *MockCourierRepo) UpdateLocation(id string, lat, lng float64) error { return nil }


func TestLogisticsUseCase_AssignCourier(t *testing.T) {
	dRepo := &MockDeliveryRepo{store: make(map[string]*logistics.Delivery)}
	cRepo := &MockCourierRepo{store: make(map[string]*logistics.Courier)}
	uc := NewLogisticsUseCase(dRepo, cRepo)

	// Setup Courier
	courier := logistics.NewCourier("Vasya", "123")
	courier.GoOnline() // Must be free
	cRepo.Save(courier)

	// Act
	err := uc.AssignCourierToDelivery(context.Background(), "order-1", courier.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Assert
	d, _ := dRepo.FindByOrderID("order-1")
	if d == nil {
		t.Fatal("delivery not created")
	}
	if d.CourierID() != courier.ID() {
		t.Error("courier not assigned to delivery")
	}
	if d.Status() != logistics.DelStatusAssigned {
		t.Error("delivery status mismatch")
	}

	c, _ := cRepo.FindByID(courier.ID())
	if c.Status() != logistics.CourierBusy {
		t.Error("courier should be busy")
	}
}
