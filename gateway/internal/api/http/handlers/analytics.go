package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	analytics_pb "github.com/versoit/diploma/services/analytics/api/proto/pb"
	"go.uber.org/zap"
)

type AnalyticsHandler struct {
	log    *zap.Logger
	client analytics_pb.KpiServiceClient
}

func NewAnalyticsHandler(log *zap.Logger, client analytics_pb.KpiServiceClient) *AnalyticsHandler {
	return &AnalyticsHandler{
		log:    log,
		client: client,
	}
}

func (h *AnalyticsHandler) GetManagerKPI(c *fiber.Ctx) error {
	id := c.Params("id")
	h.log.Info("Checking manager KPI", zap.String("manager_id", id))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.GetManagerKPI(ctx, &analytics_pb.KpiRequest{
		ManagerId: id,
	})
	if err != nil {
		h.log.Error("Analytics service error", zap.Error(err))
		return ErrorResponse(c, fiber.StatusNotFound, "manager KPI not found")
	}

	return SuccessResponse(c, resp)
}
