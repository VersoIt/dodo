package handler

import (
	"context"
	"log/slog"
	"time"

	logisticspb "github.com/versoit/diploma/be/logistics/api/proto/pb"
	"github.com/versoit/diploma/be/logistics/internal/domain"
	"github.com/versoit/diploma/be/logistics/internal/usecase"
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
	h.log.Info("Creating delivery", slog.String("order_id", req.OrderId), slog.String("order_number", req.OrderNumber))

	items := make([]domain.DeliveryItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = domain.DeliveryItem{
			ProductID: item.ProductId,
			Name:      item.ProductName,
			Quantity:  int(item.Quantity),
		}
	}

	err := h.uc.CreateDelivery(ctx, req.OrderId, req.OrderNumber, req.City, req.Street, req.House, req.Apartment, items)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create delivery: %v", err)
	}

	return &logisticspb.DeliveryResponse{OrderId: req.OrderId}, nil
}

func (h *LogisticsHandler) ListDeliveries(ctx context.Context, req *logisticspb.ListDeliveriesRequest) (*logisticspb.ListDeliveriesResponse, error) {
	deliveries, err := h.uc.ListDeliveries(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list deliveries: %v", err)
	}

	pbDeliveries := make([]*logisticspb.Delivery, len(deliveries))
	for i, d := range deliveries {
		lat, lng := d.Location()
		items := make([]*logisticspb.DeliveryItem, len(d.Items()))
		for j, it := range d.Items() {
			items[j] = &logisticspb.DeliveryItem{
				ProductId:   it.ProductID,
				ProductName: it.Name,
				Quantity:    int32(it.Quantity),
			}
		}

		pbDeliveries[i] = &logisticspb.Delivery{
			OrderId:     d.OrderID(),
			OrderNumber: d.OrderNumber(),
			CourierId:   d.CourierID(),
			Status:      d.Status().String(),
			Lat:         lat,
			Lng:         lng,
			CreatedAt:   d.CreatedAt().Format(time.RFC3339),
			City:        d.City(),
			Street:      d.Street(),
			House:       d.House(),
			Apartment:   d.Apartment(),
			Items:       items,
		}
	}

	return &logisticspb.ListDeliveriesResponse{Deliveries: pbDeliveries}, nil
}

func (h *LogisticsHandler) AssignCourier(ctx context.Context, req *logisticspb.AssignCourierRequest) (*logisticspb.DeliveryResponse, error) {
	h.log.Info("Assigning courier", slog.String("order_id", req.OrderId), slog.String("courier_id", req.CourierId))

	err := h.uc.AssignCourierToDelivery(ctx, req.OrderId, req.CourierId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "assign courier: %v", err)
	}

	return &logisticspb.DeliveryResponse{OrderId: req.OrderId, Status: "assigned"}, nil
}

func (h *LogisticsHandler) UpdateStatus(ctx context.Context, req *logisticspb.UpdateDeliveryStatusRequest) (*logisticspb.DeliveryResponse, error) {
	h.log.Info("Updating delivery status", slog.String("order_id", req.OrderId), slog.String("status", req.Status))

	var err error
	switch req.Status {
	case "delivering":
		err = h.uc.StartDelivery(ctx, req.OrderId)
	case "completed":
		err = h.uc.CompleteDelivery(ctx, req.OrderId)
	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid status: %s", req.Status)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "update status: %v", err)
	}

	return &logisticspb.DeliveryResponse{OrderId: req.OrderId, Status: req.Status}, nil
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
