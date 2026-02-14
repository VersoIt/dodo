package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	analytics_pb "github.com/versoit/diploma/services/analytics/api/proto/pb"
	"log/slog"
)

type AnalyticsHandler struct {
	log    *slog.Logger
	client analytics_pb.KpiServiceClient
}

func NewAnalyticsHandler(log *slog.Logger, client analytics_pb.KpiServiceClient) *AnalyticsHandler {
	return &AnalyticsHandler{
		log:    log,
		client: client,
	}
}

func (h *AnalyticsHandler) GetManagerKPI(c *fiber.Ctx) error {
	id := c.Params("id")
	h.log.Info("Checking manager KPI", slog.String("manager_id", id))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.GetManagerKPI(ctx, &analytics_pb.KpiRequest{
		ManagerId: id,
	})
	if err != nil {
		h.log.Error("Analytics service error", slog.Any("error", err))
		return HandleGrpcError(c, h.log, err, "failed to get analytics")
	}

	return SuccessResponse(c, resp)
}
