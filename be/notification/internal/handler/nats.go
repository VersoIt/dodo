package handler

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/versoit/diploma/be/notification/internal/usecase"
)

type OrderStatusChangedEvent struct {
	CustomerID string `json:"customer_id"`
	OrderID    string `json:"order_id"`
	Status     string `json:"status"`
}

type NatsHandler struct {
	nc  *nats.Conn
	uc  *usecase.NotificationUseCase
	log *slog.Logger
}

func NewNatsHandler(nc *nats.Conn, uc *usecase.NotificationUseCase, log *slog.Logger) *NatsHandler {
	return &NatsHandler{nc: nc, uc: uc, log: log}
}

func (h *NatsHandler) Start() error {
	_, err := h.nc.Subscribe("order.status_changed", func(msg *nats.Msg) {
		h.log.Info("Received order.status_changed event")
		var evt OrderStatusChangedEvent
		if err := json.Unmarshal(msg.Data, &evt); err != nil {
			h.log.Error("Failed to unmarshal OrderStatusChangedEvent", "error", err)
			return
		}

		title := "Order Status Update"
		body := "Your order " + evt.OrderID + " is now " + evt.Status
		err := h.uc.NotifyUser(context.Background(), evt.CustomerID, title, body)
		if err != nil {
			h.log.Error("Failed to process notification", "error", err, "order_id", evt.OrderID)
		}
	})
	return err
}
