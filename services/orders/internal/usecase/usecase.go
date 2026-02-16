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
	"github.com/versoit/diploma/services/orders/internal/domain"
)

var (
	ErrInvalidInput = errors.New("invalid input")
)

type OrderUseCase struct {
	repo           domain.OrderRepository
	catalogService domain.CatalogService
	tm             trm.Manager
	log            *slog.Logger
}

func NewOrderUseCase(repo domain.OrderRepository, catalog domain.CatalogService, tm trm.Manager, log *slog.Logger) *OrderUseCase {
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
	Address    domain.DeliveryAddress
	Items      []OrderItemInput
	PromoCode  string
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, input CreateOrderInput) (*domain.Order, error) {
	var order *domain.Order
	err := uc.tm.Do(ctx, func(ctx context.Context) error {
		order = domain.NewOrder(input.CustomerID, input.Address)

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

func (uc *OrderUseCase) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	order, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find order by id: %w", err)
	}
	return order, nil
}

func (uc *OrderUseCase) ListOrders(ctx context.Context, customerID string) ([]*domain.Order, error) {
	orders, err := uc.repo.FindByCustomerID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("find orders by customer id: %w", err)
	}
	return orders, nil
}

func (uc *OrderUseCase) ListAllOrders(ctx context.Context) ([]*domain.Order, error) {
	orders, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("find all orders: %w", err)
	}
	return orders, nil
}

func (uc *OrderUseCase) PayOrder(ctx context.Context, orderID string) error {
	return uc.tm.Do(ctx, func(ctx context.Context) error {
		order, err := uc.repo.FindByID(ctx, orderID)
		if err != nil {
			return fmt.Errorf("find order: %w", err)
		}

		if err = order.MarkPaid(); err != nil {
			return fmt.Errorf("mark paid: %w", err)
		}

		if err = uc.repo.Save(ctx, order); err != nil {
			return fmt.Errorf("save order: %w", err)
		}

		return nil
	})
}

func (uc *OrderUseCase) UpdateStatus(ctx context.Context, orderID string, statusStr string, performerID string) (*domain.Order, error) {
	status, err := domain.ParseStatus(statusStr)
	if err != nil {
		return nil, fmt.Errorf("parse status: %w", err)
	}

	var order *domain.Order
	err = uc.tm.Do(ctx, func(ctx context.Context) error {
		var err error

		order, err = uc.repo.FindByID(ctx, orderID)
		if err != nil {
			return fmt.Errorf("find order: %w", err)
		}

		var transitionErr error
		switch status {
		case domain.StatusCooking:
			transitionErr = order.SendToKitchen(performerID)
		case domain.StatusReady:
			transitionErr = order.MarkReady()
		case domain.StatusDelivering:
			transitionErr = order.ShipToDelivery(performerID)
		case domain.StatusCompleted:
			transitionErr = order.CompleteDelivery()
		default:
			return fmt.Errorf("invalid transition to status %d", status)
		}

		if transitionErr != nil {
			return fmt.Errorf("transition error: %w", transitionErr)
		}
		if err := uc.repo.Save(ctx, order); err != nil {
			return fmt.Errorf("save order: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return order, nil
}

func (uc *OrderUseCase) CreatePromoCode(ctx context.Context, code, dType string, amount float64) (*domain.PromoCode, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate uuid: %w", err)
	}

	p := domain.NewPromoCode(id.String(), strings.ToUpper(code), dType, common.NewMoney(amount), true, time.Time{})
	if err := uc.repo.SavePromo(ctx, p); err != nil {
		return nil, fmt.Errorf("save promo: %w", err)
	}

	return p, nil
}

func (uc *OrderUseCase) ListPromos(ctx context.Context) ([]*domain.PromoCode, error) {
	promos, err := uc.repo.ListPromos(ctx)
	if err != nil {
		return nil, fmt.Errorf("list promos: %w", err)
	}

	return promos, nil
}

func (uc *OrderUseCase) DeletePromo(ctx context.Context, id string) error {
	if err := uc.repo.DeletePromo(ctx, id); err != nil {
		return fmt.Errorf("delete promo: %w", err)
	}

	return nil
}

func (uc *OrderUseCase) GetPromoByCode(ctx context.Context, code string) (*domain.PromoCode, error) {
	promo, err := uc.repo.FindPromoByCode(ctx, strings.ToUpper(code))
	if err != nil {
		return nil, fmt.Errorf("find promo by code: %w", err)
	}
	return promo, nil
}

func (uc *OrderUseCase) GetAnalytics(ctx context.Context) (*domain.OrderStats, []domain.ProductStat, error) {
	kpis, err := uc.repo.GetKPIs(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("get kpis: %w", err)
	}

	top, err := uc.repo.GetTopProducts(ctx, 5)
	if err != nil {
		return kpis, nil, fmt.Errorf("get top products: %w", err)
	}

	return kpis, top, nil
}
