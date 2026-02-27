package handler

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/versoit/diploma/be/kitchen/internal/domain"
	"github.com/versoit/diploma/be/kitchen/internal/usecase"
)

type OrderItemPayload struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
}

type OrderPaidEvent struct {
	OrderID     string             `json:"order_id"`
	OrderNumber string             `json:"order_number"`
	Items       []OrderItemPayload `json:"items"`
}

type NatsHandler struct {
	nc  *nats.Conn
	uc  *usecase.KitchenUseCase
	log *slog.Logger
}

func NewNatsHandler(nc *nats.Conn, uc *usecase.KitchenUseCase, log *slog.Logger) *NatsHandler {
	return &NatsHandler{nc: nc, uc: uc, log: log}
}

func (h *NatsHandler) Start() error {
	_, err := h.nc.Subscribe("order.paid", func(msg *nats.Msg) {
		h.log.Info("Received order.paid event")
		var evt OrderPaidEvent
		if err := json.Unmarshal(msg.Data, &evt); err != nil {
			h.log.Error("Failed to unmarshal OrderPaidEvent", "error", err)
			return
		}

		items := make([]domain.KitchenItem, len(evt.Items))
		for i, item := range evt.Items {
			items[i] = domain.KitchenItem{
				ProductID: item.ProductID,
				Name:      item.ProductName,
				Quantity:  item.Quantity,
			}
		}

		_, err := h.uc.AcceptOrder(context.Background(), evt.OrderID, evt.OrderNumber, items)
		if err != nil {
			h.log.Error("Failed to accept order in kitchen", "error", err, "order_id", evt.OrderID)
		}
	})
	return err
}
