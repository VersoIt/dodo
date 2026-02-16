package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	kitchen_pb "github.com/versoit/diploma/services/kitchen/api/proto/pb"
	orders_pb "github.com/versoit/diploma/services/orders/api/proto/pb"
	"log/slog"
)

type KitchenHandler struct {
	log         *slog.Logger
	client      kitchen_pb.TicketServiceClient
	orderClient orders_pb.OrderServiceClient
}

func NewKitchenHandler(log *slog.Logger, client kitchen_pb.TicketServiceClient, orderClient orders_pb.OrderServiceClient) *KitchenHandler {
	return &KitchenHandler{
		log:         log,
		client:      client,
		orderClient: orderClient,
	}
}

func (h *KitchenHandler) ListTickets(c *fiber.Ctx) error {
	h.log.Info("Listing kitchen tickets via gRPC")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.client.ListTickets(ctx, &kitchen_pb.ListTicketsRequest{})
	if err != nil {
		h.log.Error("Kitchen service error", slog.Any("error", err))
		return HandleGrpcError(c, h.log, err, "failed to list tickets")
	}

	return SuccessResponse(c, resp)
}

func (h *KitchenHandler) UpdateStatus(c *fiber.Ctx) error {
	ticketID := c.Params("id")
	type Request struct {
		Status string `json:"status"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "invalid request")
	}

	h.log.Info("Updating kitchen ticket status via gRPC",
		slog.String("ticket_id", ticketID),
		slog.String("status", req.Status))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := h.client.UpdateTicketStatus(ctx, &kitchen_pb.UpdateTicketStatusRequest{
		TicketId: ticketID,
		Status:   req.Status,
	})
	if err != nil {
		h.log.Error("Kitchen service error", slog.Any("error", err))
		return HandleGrpcError(c, h.log, err, "failed to update status")
	}

	// Sync with Order Service
	userID := c.Locals("user_id").(string)
	_, orderErr := h.orderClient.UpdateOrderStatus(ctx, &orders_pb.UpdateOrderStatusRequest{
		OrderId:     resp.OrderId,
		Status:      req.Status,
		PerformerId: userID,
	})
	if orderErr != nil {
		h.log.Warn("Failed to sync status with Order service", slog.Any("error", orderErr))
	}

	return SuccessResponse(c, resp)
}
