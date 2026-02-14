package grpc

import (
	"context"
	"log/slog"

	"github.com/versoit/diploma/services/orders"
	orders_pb "github.com/versoit/diploma/services/orders/api/proto/pb"
	"github.com/versoit/diploma/services/orders/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrdersHandler struct {
	orders_pb.UnimplementedOrderServiceServer
	uc  *usecase.OrderUseCase
	log *slog.Logger
}

func NewOrdersHandler(uc *usecase.OrderUseCase, log *slog.Logger) *OrdersHandler {
	return &OrdersHandler{
		uc:  uc,
		log: log,
	}
}

func (h *OrdersHandler) Register(server *grpc.Server) {
	orders_pb.RegisterOrderServiceServer(server, h)
}

func (h *OrdersHandler) CreateOrder(ctx context.Context, req *orders_pb.CreateOrderRequest) (*orders_pb.OrderResponse, error) {
	h.log.Info("Creating order", slog.String("customer_id", req.CustomerId))
	
	items := make([]usecase.OrderItemInput, len(req.Items))
	for i, item := range req.Items {
		items[i] = usecase.OrderItemInput{
			ProductID: item.ProductId,
			Quantity:  int(item.Quantity),
			SizeMult:  1.0, // Default for now
		}
	}

	order, err := h.uc.CreateOrder(ctx, usecase.CreateOrderInput{
		CustomerID: req.CustomerId,
		Address: orders.DeliveryAddress{
			City:   req.Address.City,
			Street: req.Address.Street,
		},
		Items: items,
	})
	if err != nil {
		h.log.Error("Failed to create order", slog.Any("error", err))
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	return &orders_pb.OrderResponse{
		OrderId:     order.ID(),
		Status:      order.Status().String(),
		FinalPrice:  order.FinalPrice().InexactFloat64(),
		OrderNumber: order.OrderNumber(),
	}, nil
}

func (h *OrdersHandler) PayOrder(ctx context.Context, req *orders_pb.PayOrderRequest) (*orders_pb.OrderResponse, error) {
	h.log.Info("Processing payment", slog.String("order_id", req.OrderId))
	
	err := h.uc.PayOrder(ctx, req.OrderId)
	if err != nil {
		h.log.Error("Failed to pay order", slog.String("order_id", req.OrderId), slog.Any("error", err))
		return nil, status.Errorf(codes.Internal, "failed to pay order: %v", err)
	}

	return &orders_pb.OrderResponse{OrderId: req.OrderId, Status: "paid"}, nil
}

func (h *OrdersHandler) GetOrder(ctx context.Context, req *orders_pb.GetOrderRequest) (*orders_pb.OrderResponse, error) {
	order, err := h.uc.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
	}

	return &orders_pb.OrderResponse{
		OrderId:     order.ID(),
		Status:      order.Status().String(),
		FinalPrice:  order.FinalPrice().InexactFloat64(),
		OrderNumber: order.OrderNumber(),
	}, nil
}
