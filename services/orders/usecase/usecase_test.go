package usecase

import (
	"context"
	"testing"

	"github.com/versoit/diploma/services/orders"
)

type MockOrderRepo struct {
	store map[string]*orders.Order
}

func NewMockRepo() *MockOrderRepo {
	return &MockOrderRepo{
		store: make(map[string]*orders.Order),
	}
}

func (m *MockOrderRepo) Save(ctx context.Context, o *orders.Order) error {
	m.store[o.ID()] = o
	return nil
}

func (m *MockOrderRepo) FindByID(ctx context.Context, id string) (*orders.Order, error) {
	if o, ok := m.store[id]; ok {
		return o, nil
	}
	return nil, orders.ErrOrderNotFound
}

func TestOrderUseCase_CreateOrder(t *testing.T) {
	repo := NewMockRepo()
	uc := NewOrderUseCase(repo)

	input := CreateOrderInput{
		CustomerID: "cust1",
		Address: orders.DeliveryAddress{
			City:   "Moscow",
			Street: "Red Square",
		},
		Items: []OrderItemInput{
			{
				ProductID: "p1",
				Name:      "Pizza",
				Quantity:  1,
				BasePrice: 500,
				SizeMult:  1.0,
			},
		},
	}

	order, err := uc.CreateOrder(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if order.ID() == "" {
		t.Error("order ID should be generated")
	}

	savedOrder, _ := repo.FindByID(context.Background(), order.ID())
	if savedOrder == nil {
		t.Error("order should be saved in repo")
	}
	
	if savedOrder.FinalPrice() != 500 {
		t.Errorf("expected price 500, got %v", savedOrder.FinalPrice())
	}
}

func TestOrderUseCase_PayOrder(t *testing.T) {
	repo := NewMockRepo()
	uc := NewOrderUseCase(repo)

	order := orders.NewOrder("cust1", orders.DeliveryAddress{})
	repo.Save(context.Background(), order)

	err := uc.PayOrder(context.Background(), order.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updatedOrder, _ := repo.FindByID(context.Background(), order.ID())
	if updatedOrder.Status() != orders.StatusPaid {
		t.Errorf("expected status paid, got %v", updatedOrder.Status())
	}
}