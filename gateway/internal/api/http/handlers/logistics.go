package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	logistics_pb "github.com/versoit/diploma/services/logistics/api/proto/pb"
	"log/slog"
)

type LogisticsHandler struct {
	log    *slog.Logger
	client logistics_pb.DeliveryServiceClient
}

func NewLogisticsHandler(log *slog.Logger, client logistics_pb.DeliveryServiceClient) *LogisticsHandler {
	return &LogisticsHandler{
		log:    log,
		client: client,
	}
}

func (h *LogisticsHandler) AssignCourier(c *fiber.Ctx) error {
	orderID := c.Params("id")
	h.log.Info("Assigning courier (Creating Delivery) via gRPC", slog.String("order_id", orderID))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.CreateDelivery(ctx, &logistics_pb.CreateDeliveryRequest{
		OrderId: orderID,
	})
	if err != nil {
		h.log.Error("Logistics service error", slog.Any("error", err))
		return HandleGrpcError(c, h.log, err, "failed to assign courier")
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
