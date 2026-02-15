package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/google/uuid"
	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/services/orders"
)

var (
	ErrInvalidInput = errors.New("invalid input")
)

type OrderUseCase struct {
	repo           orders.OrderRepository
	catalogService orders.CatalogService
	tm             trm.Manager
	log            *slog.Logger
}

func NewOrderUseCase(repo orders.OrderRepository, catalog orders.CatalogService, tm trm.Manager, log *slog.Logger) *OrderUseCase {
	return &OrderUseCase{
		repo:           repo,
		catalogService: catalog,
		tm:             tm,
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
	PromoCode  string
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, input CreateOrderInput) (*orders.Order, error) {
	var order *orders.Order
	err := uc.tm.Do(ctx, func(ctx context.Context) error {
		order = orders.NewOrder(input.CustomerID, input.Address)

		for _, item := range input.Items {
			product, err := uc.catalogService.GetProduct(ctx, item.ProductID)
			if err != nil {
				return fmt.Errorf("failed to fetch product %s from catalog: %w", item.ProductID, err)
			}

			if err := order.AddItem(product.ID, product.Name, item.Quantity, product.BasePrice, item.SizeMult, nil); err != nil {
				return fmt.Errorf("failed to add item to order: %w", err)
			}
		}

		if input.PromoCode != "" {
			promo, err := uc.repo.FindPromoByCode(ctx, input.PromoCode)
			if err == nil && promo.IsActive() {
				discount := promo.CalculateDiscount(order.FinalPrice())
				if err := order.ApplyPromoCode(promo.Code(), discount); err != nil {
					uc.log.Warn("failed to apply promo code", slog.String("code", input.PromoCode), slog.Any("error", err))
				}
			}
		}

		if err := uc.repo.Save(ctx, order); err != nil {
			return fmt.Errorf("failed to save order: %w", err)
		}
		return nil
	})

	if err != nil {
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
	status, err := orders.ParseStatus(statusStr)
	if err != nil {
		return nil, err
	}
	order, err := uc.repo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

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
		return nil, fmt.Errorf("invalid transition")
	}

	if transitionErr != nil {
		return nil, transitionErr
	}
	if err := uc.repo.Save(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

// --- Promo Codes ---

func (uc *OrderUseCase) CreatePromoCode(ctx context.Context, code, dType string, amount float64) (*orders.PromoCode, error) {
	id, _ := uuid.NewV7()

	p := orders.NewPromoCode(id.String(), strings.ToUpper(code), dType, common.NewMoney(amount), true, time.Time{})
	if err := uc.repo.SavePromo(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
}

func (uc *OrderUseCase) ListPromos(ctx context.Context) ([]*orders.PromoCode, error) {
	return uc.repo.ListPromos(ctx)
}

func (uc *OrderUseCase) DeletePromo(ctx context.Context, id string) error {
	return uc.repo.DeletePromo(ctx, id)
}

func (uc *OrderUseCase) GetPromoByCode(ctx context.Context, code string) (*orders.PromoCode, error) {
	return uc.repo.FindPromoByCode(ctx, strings.ToUpper(code))
}

// --- Analytics ---

func (uc *OrderUseCase) GetAnalytics(ctx context.Context) (*orders.OrderStats, []orders.ProductStat, error) {
	kpis, err := uc.repo.GetKPIs(ctx)
	if err != nil {
		return nil, nil, err
	}
	top, err := uc.repo.GetTopProducts(ctx, 5)
	return kpis, top, err
}
