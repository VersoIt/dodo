package handler

import (
	"context"
	"log/slog"

	logisticspb "github.com/versoit/diploma/services/logistics/api/proto/pb"
	"github.com/versoit/diploma/services/logistics/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LogisticsHandler struct {
	logisticspb.UnimplementedDeliveryServiceServer
	uc  *usecase.LogisticsUseCase
	log *slog.Logger
}

func NewLogisticsHandler(uc *usecase.LogisticsUseCase, log *slog.Logger) *LogisticsHandler {
	return &LogisticsHandler{
		uc:  uc,
		log: log,
	}
}

func (h *LogisticsHandler) Register(server *grpc.Server) {
	logisticspb.RegisterDeliveryServiceServer(server, h)
}

func (h *LogisticsHandler) CreateDelivery(ctx context.Context, req *logisticspb.CreateDeliveryRequest) (*logisticspb.DeliveryResponse, error) {
	h.log.Info("Creating delivery", slog.String("order_id", req.OrderId))
	return &logisticspb.DeliveryResponse{OrderId: req.OrderId}, nil
}

func (h *LogisticsHandler) UpdateLocation(ctx context.Context, req *logisticspb.UpdateLocationRequest) (*logisticspb.DeliveryResponse, error) {
	h.log.Info("Updating courier location", slog.String("order_id", req.OrderId))

	err := h.uc.UpdateLocation(ctx, req.OrderId, req.Lat, req.Lng)
	if err != nil {
		h.log.Error("Failed to update location", slog.String("order_id", req.OrderId), slog.Any("error", err))
		return nil, status.Errorf(codes.Internal, "update location: %v", err)
	}
	return &logisticspb.DeliveryResponse{OrderId: req.OrderId}, nil
}
