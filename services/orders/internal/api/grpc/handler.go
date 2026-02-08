package grpc

import (
	"context"

	"github.com/versoit/diploma/services/orders"
	orders_pb "github.com/versoit/diploma/services/orders/api/proto/pb"
	"github.com/versoit/diploma/services/orders/usecase"
	"google.golang.org/grpc"
)

type OrdersHandler struct {
	orders_pb.UnimplementedOrderServiceServer
	uc *usecase.OrderUseCase
}

func NewOrdersHandler(uc *usecase.OrderUseCase) *OrdersHandler {
	return &OrdersHandler{uc: uc}
}

func (h *OrdersHandler) Register(server *grpc.Server) {
	orders_pb.RegisterOrderServiceServer(server, h)
}

func (h *OrdersHandler) CreateOrder(ctx context.Context, req *orders_pb.CreateOrderRequest) (*orders_pb.OrderResponse, error) {
	items := make([]usecase.OrderItemInput, len(req.Items))
	for i, item := range req.Items {
		items[i] = usecase.OrderItemInput{
			ProductID: item.ProductId,
			Name:      item.ProductName,
			Quantity:  int(item.Quantity),
			BasePrice: 0,
			SizeMult:  1.0,
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
		return nil, err
	}

	return &orders_pb.OrderResponse{
		OrderId:     order.ID(),
		Status:      order.Status().String(),
		FinalPrice:  float64(order.FinalPrice()),
		OrderNumber: order.OrderNumber(),
	}, nil
}

func (h *OrdersHandler) PayOrder(ctx context.Context, req *orders_pb.PayOrderRequest) (*orders_pb.OrderResponse, error) {
	err := h.uc.PayOrder(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}

	return &orders_pb.OrderResponse{OrderId: req.OrderId, Status: "paid"}, nil
}
