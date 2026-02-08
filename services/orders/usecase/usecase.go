package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/services/orders"
)

var (
	ErrInvalidInput = errors.New("invalid input data")
)

type CreateOrderInput struct {
	CustomerID string
	Address    orders.DeliveryAddress
	Items      []OrderItemInput
}

type OrderItemInput struct {
	ProductID string
	Name      string
	Quantity  int
	BasePrice common.Money
	SizeMult  float64
	Toppings  []orders.Topping
}

type OrderUseCase struct {
	repo orders.OrderRepository
}

func NewOrderUseCase(repo orders.OrderRepository) *OrderUseCase {
	return &OrderUseCase{repo: repo}
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, input CreateOrderInput) (*orders.Order, error) {
	// Валидация входных данных
	if input.CustomerID == "" {
		return nil, fmt.Errorf("%w: customer ID is required", ErrInvalidInput)
	}
	if len(input.Items) == 0 {
		return nil, fmt.Errorf("%w: order must have at least one item", ErrInvalidInput)
	}
	if input.Address.City == "" || input.Address.Street == "" {
		return nil, fmt.Errorf("%w: incomplete delivery address", ErrInvalidInput)
	}

	// Проверка контекста перед началом тяжелой операции
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	order := orders.NewOrder(input.CustomerID, input.Address)

	for _, item := range input.Items {
		if err := order.AddItem(
			item.ProductID,
			item.Name,
			item.Quantity,
			item.BasePrice,
			item.SizeMult,
			item.Toppings,
		); err != nil {
			return nil, fmt.Errorf("failed to add item %s to order: %w", item.ProductID, err)
		}
	}

	if err := uc.repo.Save(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to save order to repository: %w", err)
	}

	return order, nil
}

func (uc *OrderUseCase) PayOrder(ctx context.Context, orderID string) error {
	if orderID == "" {
		return fmt.Errorf("%w: order ID is required", ErrInvalidInput)
	}

	order, err := uc.repo.FindByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to find order %s: %w", orderID, err)
	}

	if err := order.MarkPaid(); err != nil {
		return fmt.Errorf("could not process payment for order %s: %w", orderID, err)
	}

	if err := uc.repo.Save(ctx, order); err != nil {
		return fmt.Errorf("failed to update order status after payment: %w", err)
	}

	return nil
}
