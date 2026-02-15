package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/versoit/diploma/services/orders"
)

var (
	ErrInvalidInput = errors.New("invalid input")
)

type OrderUseCase struct {
	repo           orders.OrderRepository
	catalogService orders.CatalogService
	log            *slog.Logger
}

func NewOrderUseCase(repo orders.OrderRepository, catalog orders.CatalogService, log *slog.Logger) *OrderUseCase {
	return &OrderUseCase{
		repo:           repo,
		catalogService: catalog,
		log:            log,
	}
}

type OrderItemInput struct {
	ProductID string
	Quantity  int
	SizeMult  float64
}

type CreateOrderInput struct {
	CustomerID string
	Address    orders.DeliveryAddress
	Items      []OrderItemInput
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, input CreateOrderInput) (*orders.Order, error) {
	order := orders.NewOrder(input.CustomerID, input.Address)

	for _, item := range input.Items {
		product, err := uc.catalogService.GetProduct(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch product %s from catalog: %w", item.ProductID, err)
		}

		err = order.AddItem(product.ID, product.Name, item.Quantity, product.BasePrice, item.SizeMult, nil)
		if err != nil {
			return nil, err
		}
	}

	if err := uc.repo.Save(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (uc *OrderUseCase) GetOrder(ctx context.Context, id string) (*orders.Order, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *OrderUseCase) ListOrders(ctx context.Context, customerID string) ([]*orders.Order, error) {
	return uc.repo.FindByCustomerID(ctx, customerID)
}

func (uc *OrderUseCase) ListAllOrders(ctx context.Context) ([]*orders.Order, error) {
	return uc.repo.FindAll(ctx)
}

func (uc *OrderUseCase) PayOrder(ctx context.Context, orderID string) error {
	order, err := uc.repo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if err := order.MarkPaid(); err != nil {
		return err
	}

	return uc.repo.Save(ctx, order)
}

func (uc *OrderUseCase) UpdateStatus(ctx context.Context, orderID string, statusStr string, performerID string) (*orders.Order, error) {
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

	uc.log.Info("Updating status in usecase", 
		slog.String("order_id", orderID), 
		slog.String("new_status", statusStr), 
		slog.String("performer", performerID))

	var transitionErr error
	switch status {
	case orders.StatusCooking:
		transitionErr = order.SendToKitchen(performerID)
	case orders.StatusReady:
		transitionErr = order.MarkReady()
	case orders.StatusDelivering:
		transitionErr = order.ShipToDelivery(performerID)
	case orders.StatusCompleted:
		transitionErr = order.CompleteDelivery()
	default:
		return nil, fmt.Errorf("%w: direct transition to %s not allowed", orders.ErrInvalidTransition, statusStr)
	}

	if transitionErr != nil {
		return nil, transitionErr
	}

	// CRITICAL FIX: Save the order after updating status and assignments!
	if err := uc.repo.Save(ctx, order); err != nil {
		uc.log.Error("Failed to save order after status update", slog.Any("error", err))
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	return order, nil
}
