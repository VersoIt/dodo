package handler

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/versoit/diploma/be/logistics/internal/domain"
	"github.com/versoit/diploma/be/logistics/internal/usecase"
)

type DeliveryAddress struct {
	City      string
	Street    string
	House     string
	Apartment string
	Floor     string
	Entrance  string
	Comment   string
}

type OrderItemPayload struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
}

type OrderReadyEvent struct {
	OrderID     string             `json:"order_id"`
	OrderNumber string             `json:"order_number"`
	Address     DeliveryAddress    `json:"address"`
	Items       []OrderItemPayload `json:"items"`
}

type NatsHandler struct {
	nc  *nats.Conn
	uc  *usecase.LogisticsUseCase
	log *slog.Logger
}

func NewNatsHandler(nc *nats.Conn, uc *usecase.LogisticsUseCase, log *slog.Logger) *NatsHandler {
	return &NatsHandler{nc: nc, uc: uc, log: log}
}

func (h *NatsHandler) Start() error {
	_, err := h.nc.Subscribe("order.ready", func(msg *nats.Msg) {
		h.log.Info("Received order.ready event")
		var evt OrderReadyEvent
		if err := json.Unmarshal(msg.Data, &evt); err != nil {
			h.log.Error("Failed to unmarshal OrderReadyEvent", "error", err)
			return
		}

		items := make([]domain.DeliveryItem, len(evt.Items))
		for i, item := range evt.Items {
			items[i] = domain.DeliveryItem{
				ProductID: item.ProductID,
				Name:      item.ProductName,
				Quantity:  item.Quantity,
			}
		}

		err := h.uc.CreateDelivery(
			context.Background(),
			evt.OrderID,
			evt.OrderNumber,
			evt.Address.City,
			evt.Address.Street,
			evt.Address.House,
			evt.Address.Apartment,
			items,
		)
		if err != nil {
			h.log.Error("Failed to create delivery", "error", err, "order_id", evt.OrderID)
		}
	})
	return err
}
