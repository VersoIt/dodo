package usecase

import (
	"context"
	"fmt"

	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/services/orders"
)

// CreateOrderInput - входные данные для создания заказа.
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

// OrderUseCase - оркестратор бизнес-логики заказов.
type OrderUseCase struct {
	repo orders.OrderRepository
	// Здесь могут быть клиенты других сервисов (через интерфейсы)
}

func NewOrderUseCase(repo orders.OrderRepository) *OrderUseCase {
	return &OrderUseCase{repo: repo}
}

// CreateOrder реализует сценарий создания нового заказа.
func (uc *OrderUseCase) CreateOrder(ctx context.Context, input CreateOrderInput) (*orders.Order, error) {
	// 1. Создаем агрегат заказа
	order := orders.NewOrder(input.CustomerID, input.Address)

	// 2. Добавляем позиции
	for _, item := range input.Items {
		if err := order.AddItem(
			item.ProductID,
			item.Name,
			item.Quantity,
			item.BasePrice,
			item.SizeMult,
			item.Toppings,
		); err != nil {
			return nil, fmt.Errorf("failed to add item %s: %w", item.ProductID, err)
		}
	}

	// 3. Сохраняем заказ
	if err := uc.repo.Save(order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	return order, nil
}

// PayOrder реализует сценарий оплаты заказа.
func (uc *OrderUseCase) PayOrder(ctx context.Context, orderID string) error {
	// 1. Получаем заказ из репозитория
	order, err := uc.repo.FindByID(orderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// 2. Выполняем доменное действие
	if err := order.MarkPaid(); err != nil {
		return err
	}

	// 3. Сохраняем изменения
	if err := uc.repo.Save(order); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	// TODO: Здесь должен быть вызов Event Bus для уведомления Kitchen и Analytics
	return nil
}
