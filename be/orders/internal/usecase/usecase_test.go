package usecase

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/be/orders/internal/domain"
)

type MockOrderRepo struct {
	store map[string]*domain.Order
}

func NewMockRepo() *MockOrderRepo {
	return &MockOrderRepo{
		store: make(map[string]*domain.Order),
	}
}

func (m *MockOrderRepo) Save(ctx context.Context, o *domain.Order) error {
	m.store[o.ID()] = o
	return nil
}

func (m *MockOrderRepo) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	if o, ok := m.store[id]; ok {
		return o, nil
	}
	return nil, domain.ErrOrderNotFound
}

func (m *MockOrderRepo) FindByCustomerID(ctx context.Context, cid string) ([]*domain.Order, error) { return nil, nil }
func (m *MockOrderRepo) FindAll(ctx context.Context) ([]*domain.Order, error) { return nil, nil }
func (m *MockOrderRepo) FindFiltered(ctx context.Context, filter domain.OrderFilter) ([]*domain.Order, error) { return nil, nil }
func (m *MockOrderRepo) SavePromo(ctx context.Context, p *domain.PromoCode) error { return nil }
func (m *MockOrderRepo) FindPromoByCode(ctx context.Context, code string) (*domain.PromoCode, error) { return nil, nil }
func (m *MockOrderRepo) ListPromos(ctx context.Context) ([]*domain.PromoCode, error) { return nil, nil }
func (m *MockOrderRepo) DeletePromo(ctx context.Context, id string) error { return nil }

type dummyCatalog struct{}
func (d *dummyCatalog) GetProduct(ctx context.Context, id string) (*domain.ProductInfo, error) {
	return &domain.ProductInfo{ID: id, Name: "Test", BasePrice: common.NewMoney(500)}, nil
}

type dummyKitchen struct{}
func (d *dummyKitchen) CreateTicket(ctx context.Context, orderID string, orderNumber string, items []*domain.OrderItem) error { return nil }

type dummyLogistics struct{}
func (d *dummyLogistics) CreateDelivery(ctx context.Context, orderID string, orderNumber string, address domain.DeliveryAddress, items []*domain.OrderItem) error { return nil }

type dummyTreasury struct{}
func (d *dummyTreasury) ProcessPayment(ctx context.Context, orderID string, amount common.Money) error { return nil }

type dummyNotify struct{}
func (d *dummyNotify) NotifyStatusChanged(ctx context.Context, customerID string, orderID string, status domain.OrderStatus) error { return nil }

type dummyTM struct{}
func (d *dummyTM) Do(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) }
func (d *dummyTM) DoWithSettings(ctx context.Context, _ trm.Settings, fn func(ctx context.Context) error) error { return fn(ctx) }

func TestOrderUseCase_CreateOrder(t *testing.T) {
	repo := NewMockRepo()
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	uc := NewOrderUseCase(repo, &dummyCatalog{}, &dummyKitchen{}, &dummyLogistics{}, &dummyTreasury{}, &dummyNotify{}, &dummyTM{}, log)

	input := CreateOrderInput{
		CustomerID: "cust1",
		Address: domain.DeliveryAddress{
			City:   "Moscow",
			Street: "Red Square",
			House:  "1",
		},
		Items: []OrderItemInput{
			{
				ProductID: "p1",
				Quantity:  1,
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

	savedOrder, err := repo.FindByID(context.Background(), order.ID())
	if err != nil {
		t.Fatalf("failed to find saved order: %v", err)
	}
	
	if !savedOrder.FinalPrice().Equal(common.NewMoney(500)) {
		t.Errorf("expected price 500, got %v", savedOrder.FinalPrice())
	}
}

func TestOrderUseCase_PayOrder(t *testing.T) {
	repo := NewMockRepo()
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	uc := NewOrderUseCase(repo, &dummyCatalog{}, &dummyKitchen{}, &dummyLogistics{}, &dummyTreasury{}, &dummyNotify{}, &dummyTM{}, log)

	order := domain.NewOrder("cust1", domain.DeliveryAddress{})
	if err := repo.Save(context.Background(), order); err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	err := uc.PayOrder(context.Background(), order.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updatedOrder, _ := repo.FindByID(context.Background(), order.ID())
	if updatedOrder.Status() != domain.StatusPaid {
		t.Errorf("expected status paid, got %v", updatedOrder.Status())
	}
}
