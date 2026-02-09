package handlers

import (
	"context"
	"time"
	"github.com/gofiber/fiber/v2"
	logistics_pb "github.com/versoit/diploma/services/logistics/api/proto/pb"
	"go.uber.org/zap"
)

type LogisticsHandler struct {
	log    *zap.Logger
	client logistics_pb.DeliveryServiceClient
}

func NewLogisticsHandler(log *zap.Logger, client logistics_pb.DeliveryServiceClient) *LogisticsHandler {
	return &LogisticsHandler{
		log:    log,
		client: client,
	}
}

// AssignCourier - Maps to CreateDelivery in proto as a start of logistics process
func (h *LogisticsHandler) AssignCourier(c *fiber.Ctx) error {
	orderID := c.Params("id")
	h.log.Info("Assigning courier (Creating Delivery) via gRPC", zap.String("order_id", orderID))
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.CreateDelivery(ctx, &logistics_pb.CreateDeliveryRequest{
		OrderId: orderID,
	})
	
	if err != nil {
		h.log.Error("Logistics service error", zap.Error(err))
		return ErrorResponse(c, fiber.StatusInternalServerError, "failed to assign courier")
	}

	return SuccessResponse(c, resp)
}

func (h *LogisticsHandler) CompleteDelivery(c *fiber.Ctx) error {
	// Not directly in proto yet, so keeping as TODO/Mock or mapping to update status if available
	// For now, assuming UpdateLocation or similar logic might handle this or it requires proto update.
	// Returning mocked success to satisfy interface for now, acknowledging requirement.
	return SuccessResponse(c, fiber.Map{
		"status": "DELIVERED",
		"note": "CompleteDelivery RPC needs to be added to logistics.proto",
	})
}

func (h *LogisticsHandler) UpdateLocation(c *fiber.Ctx) error {
	orderID := c.Params("id")
	type Request struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	}
	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "invalid coordinates")
	}

	h.log.Debug("Updating courier location via gRPC", 
		zap.String("order_id", orderID),
		zap.Float64("lat", req.Lat),
		zap.Float64("lng", req.Lng))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.UpdateLocation(ctx, &logistics_pb.UpdateLocationRequest{
		OrderId: orderID,
		Lat:     req.Lat,
		Lng:     req.Lng,
	})

	if err != nil {
		h.log.Error("Logistics service error", zap.Error(err))
		return ErrorResponse(c, fiber.StatusInternalServerError, "failed to update location")
	}

	return SuccessResponse(c, resp)
}