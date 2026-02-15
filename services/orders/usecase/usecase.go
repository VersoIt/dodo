package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

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
	Quantity  int
	SizeMult  float64
	Toppings  []orders.Topping
}

type OrderUseCase struct {
	repo      orders.OrderRepository
	catalog   orders.CatalogService
	analytics orders.AnalyticsService
	log       *slog.Logger
}

func NewOrderUseCase(repo orders.OrderRepository, catalog orders.CatalogService, analytics orders.AnalyticsService, log *slog.Logger) *OrderUseCase {
	return &OrderUseCase{
		repo:      repo,
		catalog:   catalog,
		analytics: analytics,
		log:       log,
	}
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
		// Fetch actual product info from catalog
		prod, err := uc.catalog.GetProduct(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch product %s from catalog: %w", item.ProductID, err)
		}

		if err := order.AddItem(
			prod.ID,
			prod.Name,
			item.Quantity,
			prod.BasePrice,
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

	// Report to analytics (Best effort or should it be transactional?
	// In a real senior scenario we'd use Outbox pattern, but here we'll do a direct call for now)
	if err := uc.analytics.ReportSale(ctx, "019c53ee-74af-72b9-afe6-103c1466ae0e", order.FinalPrice().InexactFloat64()); err != nil {
		// We log it but don't fail the payment as it's already saved in DB
		uc.log.Error("failed to report sale to analytics", slog.Any("error", err), slog.String("order_id", orderID))
	}

	return nil
}

func (uc *OrderUseCase) GetOrder(ctx context.Context, id string) (*orders.Order, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: order ID is required", ErrInvalidInput)
	}

	order, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find order %s: %w", id, err)
	}

	return order, nil
}

func (uc *OrderUseCase) ListOrders(ctx context.Context, customerID string) ([]*orders.Order, error) {
	if customerID == "" {
		return nil, fmt.Errorf("%w: customer ID is required", ErrInvalidInput)
	}

	return uc.repo.FindByCustomerID(ctx, customerID)
}

func (uc *OrderUseCase) ListAllOrders(ctx context.Context) ([]*orders.Order, error) {
	return uc.repo.FindAll(ctx)
}

func (uc *OrderUseCase) UpdateStatus(ctx context.Context, orderID string, statusStr string) (*orders.Order, error) {
	if orderID == "" {
		return nil, fmt.Errorf("%w: order ID is required", ErrInvalidInput)
	}

	status, err := orders.ParseStatus(statusStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	order, err := uc.repo.FindByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find order %s: %w", orderID, err)
	}

	switch status {
	case orders.StatusCooking:
		err = order.SendToKitchen()
	case orders.StatusReady:
		err = order.MarkReady()
	case orders.StatusDelivering:
		err = order.ShipToDelivery()
	case orders.StatusCompleted:
		err = order.CompleteDelivery()
	default:
		return nil, fmt.Errorf("%w: direct transition to %s not allowed", orders.ErrInvalidTransition, statusStr)
	}

	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to save order %s: %w", orderID, err)
	}

	return order, nil
}
