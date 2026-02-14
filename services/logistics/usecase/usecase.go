package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/versoit/diploma/services/logistics"
)

var (
	ErrInvalidInput = errors.New("invalid input data")
)

type LogisticsUseCase struct {
	deliveryRepo logistics.DeliveryRepository
	courierRepo  logistics.CourierRepository
	log          *slog.Logger
}

func NewLogisticsUseCase(dr logistics.DeliveryRepository, cr logistics.CourierRepository, log *slog.Logger) *LogisticsUseCase {
	return &LogisticsUseCase{
		deliveryRepo: dr,
		courierRepo:  cr,
		log:          log,
	}
}

func (uc *LogisticsUseCase) AssignCourierToDelivery(ctx context.Context, orderID string, courierID string) error {
	if orderID == "" || courierID == "" {
		return fmt.Errorf("%w: order ID and courier ID are required", ErrInvalidInput)
	}

	uc.log.Info("assigning courier", slog.String("order_id", orderID), slog.String("courier_id", courierID))

	delivery, err := uc.deliveryRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		delivery = logistics.NewDelivery(orderID)
	}

	courier, err := uc.courierRepo.FindByID(ctx, courierID)
	if err != nil {
		return fmt.Errorf("failed to locate courier %s: %w", courierID, err)
	}

	if err := courier.TakeOrder(); err != nil {
		return fmt.Errorf("courier %s cannot take order: %w", courierID, err)
	}

	if err := delivery.AssignCourier(courier.ID()); err != nil {
		return fmt.Errorf("delivery assignment failed: %w", err)
	}

	if err := uc.courierRepo.Save(ctx, courier); err != nil {
		return fmt.Errorf("failed to update courier status: %w", err)
	}
	if err := uc.deliveryRepo.Save(ctx, delivery); err != nil {
		return fmt.Errorf("failed to save delivery assignment: %w", err)
	}

	return nil
}

func (uc *LogisticsUseCase) UpdateLocation(ctx context.Context, orderID string, lat, lng float64) error {
	if orderID == "" {
		return fmt.Errorf("%w: order ID is required", ErrInvalidInput)
	}

	delivery, err := uc.deliveryRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("delivery process not found for order %s: %w", orderID, err)
	}

	if err := delivery.UpdateLocation(lat, lng); err != nil {
		return fmt.Errorf("failed to update delivery coordinates: %w", err)
	}

	if err := uc.deliveryRepo.Save(ctx, delivery); err != nil {
		return fmt.Errorf("failed to persist location update: %w", err)
	}

	return nil
}
