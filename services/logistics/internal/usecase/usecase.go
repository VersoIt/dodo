package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/versoit/diploma/services/logistics/internal/domain"
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

func (uc *LogisticsUseCase) AssignCourierToDelivery(ctx context.Context, orderID string, courierID string) error {
	if orderID == "" || courierID == "" {
		return fmt.Errorf("%w: order ID and courier ID are required", ErrInvalidInput)
	}

	uc.log.Info("assigning courier", slog.String("order_id", orderID), slog.String("courier_id", courierID))

	return uc.tm.Do(ctx, func(ctx context.Context) error {
		delivery, err := uc.deliveryRepo.FindByOrderID(ctx, orderID)
		if err != nil {
			delivery = domain.NewDelivery(orderID)
		}

		courier, err := uc.courierRepo.FindByID(ctx, courierID)
		if err != nil {
			return fmt.Errorf("locate courier %s: %w", courierID, err)
		}

		if err := courier.TakeOrder(); err != nil {
			return fmt.Errorf("courier %s take order: %w", courierID, err)
		}

		if err := delivery.AssignCourier(courier.ID()); err != nil {
			return fmt.Errorf("delivery assignment: %w", err)
		}

		if err := uc.courierRepo.Save(ctx, courier); err != nil {
			return fmt.Errorf("save courier: %w", err)
		}
		if err := uc.deliveryRepo.Save(ctx, delivery); err != nil {
			return fmt.Errorf("save delivery: %w", err)
		}
		return nil
	})
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
