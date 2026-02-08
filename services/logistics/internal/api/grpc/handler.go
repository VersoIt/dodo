package grpc

import (
	"context"

	logistics_pb "github.com/versoit/diploma/services/logistics/api/proto/pb"
	"github.com/versoit/diploma/services/logistics/usecase"
	"google.golang.org/grpc"
)

type LogisticsHandler struct {
	logistics_pb.UnimplementedDeliveryServiceServer
	uc *usecase.LogisticsUseCase
}

func NewLogisticsHandler(uc *usecase.LogisticsUseCase) *LogisticsHandler {
	return &LogisticsHandler{uc: uc}
}

func (h *LogisticsHandler) Register(server *grpc.Server) {
	logistics_pb.RegisterDeliveryServiceServer(server, h)
}

func (h *LogisticsHandler) CreateDelivery(ctx context.Context, req *logistics_pb.CreateDeliveryRequest) (*logistics_pb.DeliveryResponse, error) {
	return &logistics_pb.DeliveryResponse{OrderId: req.OrderId}, nil
}

func (h *LogisticsHandler) UpdateLocation(ctx context.Context, req *logistics_pb.UpdateLocationRequest) (*logistics_pb.DeliveryResponse, error) {
	err := h.uc.UpdateLocation(ctx, req.OrderId, req.Lat, req.Lng)
	if err != nil {
		return nil, err
	}
	return &logistics_pb.DeliveryResponse{OrderId: req.OrderId}, nil
}
