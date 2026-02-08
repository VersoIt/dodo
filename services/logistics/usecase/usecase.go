package usecase

import (
	"context"
	"fmt"

	"github.com/versoit/diploma/services/logistics"
)

type LogisticsUseCase struct {
	deliveryRepo logistics.DeliveryRepository
	courierRepo  logistics.CourierRepository
}

func NewLogisticsUseCase(dr logistics.DeliveryRepository, cr logistics.CourierRepository) *LogisticsUseCase {
	return &LogisticsUseCase{
		deliveryRepo: dr,
		courierRepo:  cr,
	}
}

// AssignCourierToDelivery подбирает курьера и назначает его на доставку.
func (uc *LogisticsUseCase) AssignCourierToDelivery(ctx context.Context, orderID string, courierID string) error {
	delivery, err := uc.deliveryRepo.FindByOrderID(orderID)
	if err != nil {
		// Если доставки еще нет, создаем её
		delivery = logistics.NewDelivery(orderID)
	}

	courier, err := uc.courierRepo.FindByID(courierID)
	if err != nil {
		return fmt.Errorf("courier not found: %w", err)
	}

	// Доменная логика курьера
	if err := courier.TakeOrder(); err != nil {
		return err
	}

	// Доменная логика доставки
	if err := delivery.AssignCourier(courier.ID()); err != nil {
		return err
	}

	// Сохраняем обоих
	if err := uc.courierRepo.Save(courier); err != nil {
		return err
	}
	return uc.deliveryRepo.Save(delivery)
}

// UpdateLocation обновляет координаты в процессе доставки.
func (uc *LogisticsUseCase) UpdateLocation(ctx context.Context, orderID string, lat, lng float64) error {
	delivery, err := uc.deliveryRepo.FindByOrderID(orderID)
	if err != nil {
		return err
	}

	if err := delivery.UpdateLocation(lat, lng); err != nil {
		return err
	}

	return uc.deliveryRepo.Save(delivery)
}
