package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/versoit/diploma/be/logistics/internal/domain"
)

var (
	ErrInvalidInput = errors.New("invalid input data")
)

type LogisticsUseCase struct {
	deliveryRepo domain.DeliveryRepository
	courierRepo  domain.CourierRepository
	tm           trm.Manager
	log          *slog.Logger
}

func NewLogisticsUseCase(dr domain.DeliveryRepository, cr domain.CourierRepository, tm trm.Manager, log *slog.Logger) *LogisticsUseCase {
	return &LogisticsUseCase{
		deliveryRepo: dr,
		courierRepo:  cr,
		tm:           tm,
		log:          log,
	}
}

func (uc *LogisticsUseCase) CreateDelivery(ctx context.Context, orderID, orderNumber, city, street, house, apartment string, items []domain.DeliveryItem) error {
	if orderID == "" {
		return fmt.Errorf("%w: order ID is required", ErrInvalidInput)
	}

	return uc.tm.Do(ctx, func(ctx context.Context) error {
		delivery := domain.NewDelivery(orderID, orderNumber, city, street, house, apartment, items)

		uc.log.Info("creating delivery", slog.String("order_id", orderID), slog.String("order_number", orderNumber))

		if err := uc.deliveryRepo.Save(ctx, delivery); err != nil {
			return fmt.Errorf("failed to create delivery: %w", err)
		}
		return nil
	})
}

func (uc *LogisticsUseCase) AssignCourierToDelivery(ctx context.Context, orderID string, courierID string) error {
	return uc.tm.Do(ctx, func(ctx context.Context) error {
		// Ensure courier exists (simple sync for prototype)
		if _, err := uc.courierRepo.FindByID(ctx, courierID); err != nil {
			stubCourier := domain.ReconstructCourier(courierID, "Unknown Courier", "N/A", domain.CourierFree, 0, 0)
			if err := uc.courierRepo.Save(ctx, stubCourier); err != nil {
				return fmt.Errorf("failed to sync courier: %w", err)
			}
		}

		delivery, err := uc.deliveryRepo.FindByOrderID(ctx, orderID)
		if err != nil {
			return fmt.Errorf("find delivery %s: %w", orderID, err)
		}

		if err := delivery.AssignCourier(courierID); err != nil {
			return fmt.Errorf("assign courier: %w", err)
		}

		if err := uc.deliveryRepo.Save(ctx, delivery); err != nil {
			return fmt.Errorf("persist delivery: %w", err)
		}
		return nil
	})
}

func (uc *LogisticsUseCase) StartDelivery(ctx context.Context, orderID string) error {
	return uc.tm.Do(ctx, func(ctx context.Context) error {
		delivery, err := uc.deliveryRepo.FindByOrderID(ctx, orderID)
		if err != nil {
			return fmt.Errorf("find delivery %s: %w", orderID, err)
		}

		if err := delivery.Pickup(); err != nil {
			return fmt.Errorf("pickup delivery: %w", err)
		}

		if err := uc.deliveryRepo.Save(ctx, delivery); err != nil {
			return fmt.Errorf("persist delivery: %w", err)
		}
		return nil
	})
}

func (uc *LogisticsUseCase) CompleteDelivery(ctx context.Context, orderID string) error {
	return uc.tm.Do(ctx, func(ctx context.Context) error {
		delivery, err := uc.deliveryRepo.FindByOrderID(ctx, orderID)
		if err != nil {
			return fmt.Errorf("find delivery %s: %w", orderID, err)
		}

		if err := delivery.Complete(); err != nil {
			return fmt.Errorf("complete delivery: %w", err)
		}

		if err := uc.deliveryRepo.Save(ctx, delivery); err != nil {
			return fmt.Errorf("persist delivery: %w", err)
		}
		return nil
	})
}

func (uc *LogisticsUseCase) ListDeliveries(ctx context.Context) ([]*domain.Delivery, error) {
	return uc.deliveryRepo.FindAll(ctx)
}

func (uc *LogisticsUseCase) UpdateLocation(ctx context.Context, orderID string, lat, lng float64) error {
	if orderID == "" {
		return fmt.Errorf("%w: order ID is required", ErrInvalidInput)
	}

	coords, err := domain.NewCoordinates(lat, lng)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	return uc.tm.Do(ctx, func(ctx context.Context) error {
		delivery, err := uc.deliveryRepo.FindByOrderID(ctx, orderID)
		if err != nil {
			return fmt.Errorf("find delivery %s: %w", orderID, err)
		}

		delivery.UpdateLocation(coords)

		if err := uc.deliveryRepo.Save(ctx, delivery); err != nil {
			return fmt.Errorf("persist location: %w", err)
		}
		return nil
	})
}
