package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	logistics_pb "github.com/versoit/diploma/services/logistics/api/proto/pb"
	orders_pb "github.com/versoit/diploma/services/orders/api/proto/pb"
	"log/slog"
)

type LogisticsHandler struct {
	log         *slog.Logger
	client      logistics_pb.DeliveryServiceClient
	orderClient orders_pb.OrderServiceClient
}

func NewLogisticsHandler(log *slog.Logger, client logistics_pb.DeliveryServiceClient, orderClient orders_pb.OrderServiceClient) *LogisticsHandler {
	return &LogisticsHandler{
		log:         log,
		client:      client,
		orderClient: orderClient,
	}
}

func (h *LogisticsHandler) ListDeliveries(c *fiber.Ctx) error {
	h.log.Info("Listing deliveries via gRPC")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.client.ListDeliveries(ctx, &logistics_pb.ListDeliveriesRequest{})
	if err != nil {
		h.log.Error("Logistics service error", slog.Any("error", err))
		return HandleGrpcError(c, h.log, err, "failed to list deliveries")
	}

	return SuccessResponse(c, resp)
}

func (h *LogisticsHandler) AssignCourier(c *fiber.Ctx) error {
	orderID := c.Params("id")
	type Request struct {
		CourierID string `json:"courier_id"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "invalid request")
	}

	h.log.Info("Assigning courier via gRPC", slog.String("order_id", orderID), slog.String("courier_id", req.CourierID))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.AssignCourier(ctx, &logistics_pb.AssignCourierRequest{
		OrderId:   orderID,
		CourierId: req.CourierID,
	})
	if err != nil {
		h.log.Error("Logistics service error", slog.Any("error", err))
		return HandleGrpcError(c, h.log, err, "failed to assign courier")
	}

	return SuccessResponse(c, resp)
}

func (h *LogisticsHandler) UpdateStatus(c *fiber.Ctx) error {
	orderID := c.Params("id")
	type Request struct {
		Status string `json:"status"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "invalid request")
	}

	h.log.Info("Updating delivery status via gRPC", slog.String("order_id", orderID), slog.String("status", req.Status))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := h.client.UpdateStatus(ctx, &logistics_pb.UpdateDeliveryStatusRequest{
		OrderId: orderID,
		Status:  req.Status,
	})
	if err != nil {
		h.log.Error("Logistics service error", slog.Any("error", err))
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

func (h *LogisticsHandler) UpdateLocation(c *fiber.Ctx) error {
	orderID := c.Params("id")
	type Request struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "invalid request")
	}

	h.log.Info("Updating courier location via gRPC",
		slog.String("order_id", orderID),
		slog.Float64("lat", req.Lat),
		slog.Float64("lng", req.Lng))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.UpdateLocation(ctx, &logistics_pb.UpdateLocationRequest{
		OrderId: orderID,
		Lat:     req.Lat,
		Lng:     req.Lng,
	})
	if err != nil {
		h.log.Error("Logistics service error", slog.Any("error", err))
		return HandleGrpcError(c, h.log, err, "failed to update location")
	}

	return SuccessResponse(c, resp)
}
