package grpc

import (
	"context"
	"log/slog"

	logistics_pb "github.com/versoit/diploma/services/logistics/api/proto/pb"
	"github.com/versoit/diploma/services/logistics/usecase"
	"google.golang.org/grpc"
)

type LogisticsHandler struct {
	logistics_pb.UnimplementedDeliveryServiceServer
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
	logistics_pb.RegisterDeliveryServiceServer(server, h)
}

func (h *LogisticsHandler) CreateDelivery(ctx context.Context, req *logistics_pb.CreateDeliveryRequest) (*logistics_pb.DeliveryResponse, error) {
	h.log.Info("Creating delivery", slog.String("order_id", req.OrderId))
	return &logistics_pb.DeliveryResponse{OrderId: req.OrderId}, nil
}

func (h *LogisticsHandler) UpdateLocation(ctx context.Context, req *logistics_pb.UpdateLocationRequest) (*logistics_pb.DeliveryResponse, error) {
	h.log.Info("Updating courier location", slog.String("order_id", req.OrderId))
	
	err := h.uc.UpdateLocation(ctx, req.OrderId, req.Lat, req.Lng)
	if err != nil {
		h.log.Error("Failed to update location", slog.String("order_id", req.OrderId), slog.Any("error", err))
		return nil, err
	}
	return &logistics_pb.DeliveryResponse{OrderId: req.OrderId}, nil
}
