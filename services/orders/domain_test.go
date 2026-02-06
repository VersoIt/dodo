package orders

import (
	"diploma/pkg/common"
	"testing"
)

func TestOrder_AddItem_CalculatesPriceCorrectly(t *testing.T) {
	// Arrange
	addr := DeliveryAddress{City: "Test City", Street: "Main St"}
	order := NewOrder("cust-123", addr)
	
	basePrice := common.Money(100.0)
	sizeMult := 1.2 // +20%
	qty := 2
	toppings := []Topping{
		{Name: "Cheese", Price: 10.0},
		{Name: "Sauce", Price: 5.0},
	}
	// Unit price calculation:
	// Sized Price = 100 * 1.2 = 120
	// Toppings = 10 + 5 = 15
	// Unit Price = 120 + 15 = 135
	// Total Item Price = 135 * 2 = 270

	// Act
	err := order.AddItem("prod-1", "Pizza", qty, basePrice, sizeMult, toppings)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedPrice := common.Money(270.0)
	if order.FinalPrice() != expectedPrice {
		t.Errorf("expected final price %v, got %v", expectedPrice, order.FinalPrice())
	}
}

func TestOrder_ApplyPromoCode(t *testing.T) {
	order := NewOrder("cust-1", DeliveryAddress{})
	_ = order.AddItem("p1", "Item", 1, 100, 1.0, nil) // Total 100

	err := order.ApplyPromoCode("PROMO10", 10.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := common.Money(90.0) // 100 - 10
	if order.FinalPrice() != expected {
		t.Errorf("expected %v, got %v", expected, order.FinalPrice())
	}
}

func TestOrder_StateTransitions(t *testing.T) {
	order := NewOrder("c1", DeliveryAddress{})
	
	// Created -> Paid
	if err := order.MarkPaid(); err != nil {
		t.Errorf("failed to mark paid: %v", err)
	}
	if order.status != StatusPaid {
		t.Errorf("expected status Paid, got %v", order.status)
	}

	// Paid -> Cooking
	if err := order.SendToKitchen(); err != nil {
		t.Errorf("failed to send to kitchen: %v", err)
	}
	
	// Cooking -> Ready
	if err := order.MarkReady(); err != nil {
		t.Errorf("failed to mark ready: %v", err)
	}

	// Ready -> Delivering
	if err := order.ShipToDelivery(); err != nil {
		t.Errorf("failed to ship: %v", err)
	}

	// Delivering -> Completed
	if err := order.CompleteDelivery(); err != nil {
		t.Errorf("failed to complete: %v", err)
	}

	if order.status != StatusCompleted {
		t.Errorf("expected status Completed, got %v", order.status)
	}
}

func TestOrder_CannotAddItem_WhenLocked(t *testing.T) {
	order := NewOrder("c1", DeliveryAddress{})
	_ = order.MarkPaid() // Lock order

	err := order.AddItem("p1", "Item", 1, 100, 1, nil)
	if err != ErrOrderLocked {
		t.Errorf("expected ErrOrderLocked, got %v", err)
	}
}
