package grpc

import (
	"context"

	"github.com/versoit/diploma/services/logistics/api/proto/pb"
	"github.com/versoit/diploma/services/logistics/usecase"
	"google.golang.org/grpc"
)

type LogisticsHandler struct {
	pb.UnimplementedDeliveryServiceServer
	uc *usecase.LogisticsUseCase
}

func NewLogisticsHandler(uc *usecase.LogisticsUseCase) *LogisticsHandler {
	return &LogisticsHandler{uc: uc}
}

func (h *LogisticsHandler) Register(server *grpc.Server) {
	pb.RegisterDeliveryServiceServer(server, h)
}

func (h *LogisticsHandler) CreateDelivery(ctx context.Context, req *pb.CreateDeliveryRequest) (*pb.DeliveryResponse, error) {
	// Logic to start delivery
	return &pb.DeliveryResponse{OrderId: req.OrderId}, nil
}

func (h *LogisticsHandler) UpdateLocation(ctx context.Context, req *pb.UpdateLocationRequest) (*pb.DeliveryResponse, error) {
	err := h.uc.UpdateLocation(ctx, req.OrderId, req.Lat, req.Lng)
	if err != nil {
		return nil, err
	}
	return &pb.DeliveryResponse{OrderId: req.OrderId}, nil
}
