package domain

import (
	"github.com/versoit/diploma/pkg/common"
	"testing"
)

func TestOrder_AddItem_CalculatesPriceCorrectly(t *testing.T) {
	// Arrange
	addr := DeliveryAddress{City: "Test City", Street: "Main St"}
	order := NewOrder("cust-123", addr)

	basePrice := common.NewMoney(100.0)
	sizeMult := 1.2 // +20%
	qty, _ := NewQuantity(2)
	toppings := []Topping{
		{Name: "Cheese", Price: common.NewMoney(10.0)},
		{Name: "Sauce", Price: common.NewMoney(5.0)},
	}

	// Act
	err := order.AddItem("prod-1", "Pizza", qty, basePrice, sizeMult, toppings)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedPrice := common.NewMoney(270.0)
	if !order.FinalPrice().Equal(expectedPrice) {
		t.Errorf("expected final price %v, got %v", expectedPrice, order.FinalPrice())
	}
}

func TestOrder_ApplyPromoCode(t *testing.T) {
	order := NewOrder("cust-1", DeliveryAddress{})
	qty, _ := NewQuantity(1)
	_ = order.AddItem("p1", "Item", qty, common.NewMoney(100), 1.0, nil) // Total 100

	err := order.ApplyPromoCode("PROMO10", common.NewMoney(10.0))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := common.NewMoney(90.0) // 100 - 10
	if !order.FinalPrice().Equal(expected) {
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
	if err := order.SendToKitchen("chef-1"); err != nil {
		t.Errorf("failed to send to kitchen: %v", err)
	}

	// Cooking -> Ready
	if err := order.MarkReady(); err != nil {
		t.Errorf("failed to mark ready: %v", err)
	}

	// Ready -> Delivering
	if err := order.ShipToDelivery("courier-1"); err != nil {
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

	qty, _ := NewQuantity(1)
	err := order.AddItem("p1", "Item", qty, common.NewMoney(100), 1, nil)
	if err != ErrOrderLocked {
		t.Errorf("expected ErrOrderLocked, got %v", err)
	}
}
