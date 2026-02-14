package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	kitchen_pb "github.com/versoit/diploma/services/kitchen/api/proto/pb"
	"log/slog"
)

type KitchenHandler struct {
	log    *slog.Logger
	client kitchen_pb.TicketServiceClient
}

func NewKitchenHandler(log *slog.Logger, client kitchen_pb.TicketServiceClient) *KitchenHandler {
	return &KitchenHandler{
		log:    log,
		client: client,
	}
}

func (h *KitchenHandler) UpdateStatus(c *fiber.Ctx) error {
	ticketID := c.Params("id")
	type Request struct {
		Status int32 `json:"status"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "invalid request")
	}

	h.log.Info("Updating kitchen ticket status via gRPC",
		slog.String("ticket_id", ticketID),
		slog.Int("status", int(req.Status)))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.UpdateTicketStatus(ctx, &kitchen_pb.UpdateTicketStatusRequest{
		TicketId: ticketID,
		Status:   req.Status,
	})
	if err != nil {
		h.log.Error("Kitchen service error", slog.Any("error", err))
		return HandleGrpcError(c, h.log, err, "failed to update status")
	}

	return SuccessResponse(c, resp)
}
