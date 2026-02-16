package handler

import (
	"context"
	"log/slog"

	orderspb "github.com/versoit/diploma/services/orders/api/proto/pb"
	"github.com/versoit/diploma/services/orders/internal/domain"
	"github.com/versoit/diploma/services/orders/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrdersHandler struct {
	orderspb.UnimplementedOrderServiceServer
	uc  *usecase.OrderUseCase
	log *slog.Logger
}

func NewOrdersHandler(uc *usecase.OrderUseCase, log *slog.Logger) *OrdersHandler {
	return &OrdersHandler{uc: uc, log: log}
}

func (h *OrdersHandler) Register(server *grpc.Server) {
	orderspb.RegisterOrderServiceServer(server, h)
}

func (h *OrdersHandler) CreateOrder(ctx context.Context, req *orderspb.CreateOrderRequest) (*orderspb.OrderResponse, error) {
	h.log.Info("Creating new order", slog.String("customer_id", req.CustomerId))
	items := make([]usecase.OrderItemInput, len(req.Items))
	for i, item := range req.Items {
		items[i] = usecase.OrderItemInput{ProductID: item.ProductId, Quantity: int(item.Quantity), SizeMult: 1.0}
	}

	addr := domain.DeliveryAddress{
		City:      req.Address.City,
		Street:    req.Address.Street,
		House:     req.Address.House,
		Apartment: req.Address.Apartment,
		Floor:     req.Address.Floor,
		Entrance:  req.Address.Entrance,
		Comment:   req.Address.Comment,
	}

	order, err := h.uc.CreateOrder(ctx, usecase.CreateOrderInput{
		CustomerID: req.CustomerId,
		Address:    addr,
		Items:      items,
		PromoCode:  req.PromoCode,
	})
	if err != nil {
		h.log.Error("failed to create order", slog.Any("error", err))
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}
	return h.mapOrder(order), nil
}

func (h *OrdersHandler) GetOrder(ctx context.Context, req *orderspb.GetOrderRequest) (*orderspb.OrderResponse, error) {
	order, err := h.uc.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found")
	}
	return h.mapOrder(order), nil
}

func (h *OrdersHandler) ListOrders(ctx context.Context, req *orderspb.ListOrdersRequest) (*orderspb.ListOrdersResponse, error) {
	orderList, err := h.uc.ListOrders(ctx, req.CustomerId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list orders")
	}
	pbOrders := make([]*orderspb.OrderResponse, len(orderList))
	for i, o := range orderList {
		pbOrders[i] = h.mapOrder(o)
	}
	return &orderspb.ListOrdersResponse{Orders: pbOrders}, nil
}

func (h *OrdersHandler) ListAllOrders(ctx context.Context, _ *orderspb.ListAllOrdersRequest) (*orderspb.ListOrdersResponse, error) {
	orderList, err := h.uc.ListAllOrders(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list all orders")
	}
	pbOrders := make([]*orderspb.OrderResponse, len(orderList))
	for i, o := range orderList {
		pbOrders[i] = h.mapOrder(o)
	}
	return &orderspb.ListOrdersResponse{Orders: pbOrders}, nil
}

func (h *OrdersHandler) UpdateOrderStatus(ctx context.Context, req *orderspb.UpdateOrderStatusRequest) (*orderspb.OrderResponse, error) {
	order, err := h.uc.UpdateStatus(ctx, req.OrderId, req.Status, req.PerformerId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update status: %v", err)
	}
	return h.mapOrder(order), nil
}

func (h *OrdersHandler) PayOrder(ctx context.Context, req *orderspb.PayOrderRequest) (*orderspb.OrderResponse, error) {
	err := h.uc.PayOrder(ctx, req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "payment failed: %v", err)
	}
	return &orderspb.OrderResponse{OrderId: req.OrderId, Status: "paid"}, nil
}

// --- Promo Codes ---

func (h *OrdersHandler) CreatePromoCode(ctx context.Context, req *orderspb.CreatePromoCodeRequest) (*orderspb.PromoCodeResponse, error) {
	p, err := h.uc.CreatePromoCode(ctx, req.Code, req.Type, req.Amount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create promo")
	}
	return &orderspb.PromoCodeResponse{Id: p.ID(), Code: p.Code(), Type: p.DiscountType(), Amount: p.DiscountAmount().InexactFloat64(), Active: p.IsActive()}, nil
}

func (h *OrdersHandler) ListPromos(ctx context.Context, _ *orderspb.ListPromosRequest) (*orderspb.ListPromosResponse, error) {
	promos, err := h.uc.ListPromos(ctx)
	if err != nil {
		h.log.Error("failed to list promos", slog.Any("error", err))
		return nil, status.Errorf(codes.Internal, "failed to list promos: %v", err)
	}
	pbPromos := make([]*orderspb.PromoCodeResponse, len(promos))
	for i, p := range promos {
		pbPromos[i] = &orderspb.PromoCodeResponse{Id: p.ID(), Code: p.Code(), Type: p.DiscountType(), Amount: p.DiscountAmount().InexactFloat64(), Active: p.IsActive()}
	}
	return &orderspb.ListPromosResponse{Promos: pbPromos}, nil
}

func (h *OrdersHandler) DeletePromo(ctx context.Context, req *orderspb.DeletePromoRequest) (*orderspb.PromoDeleteResponse, error) {
	if err := h.uc.DeletePromo(ctx, req.Id); err != nil {
		h.log.Error("failed to delete promo", slog.Any("error", err))
		return nil, status.Errorf(codes.Internal, "failed to delete promo: %v", err)
	}
	return &orderspb.PromoDeleteResponse{Success: true}, nil
}

func (h *OrdersHandler) CheckPromoCode(ctx context.Context, req *orderspb.CheckPromoCodeRequest) (*orderspb.PromoCodeResponse, error) {
	p, err := h.uc.GetPromoByCode(ctx, req.Code)
	if err != nil {
		h.log.Warn("promo code check failed", slog.String("code", req.Code), slog.Any("error", err))
		return nil, status.Errorf(codes.NotFound, "promo not found")
	}
	return &orderspb.PromoCodeResponse{Id: p.ID(), Code: p.Code(), Type: p.DiscountType(), Amount: p.DiscountAmount().InexactFloat64(), Active: p.IsActive()}, nil
}

func (h *OrdersHandler) mapOrder(o *domain.Order) *orderspb.OrderResponse {
	addr := o.Address()
	pbItems := make([]*orderspb.OrderItem, len(o.Items()))
	for i, item := range o.Items() {
		pbItems[i] = &orderspb.OrderItem{ProductId: item.ProductID(), ProductName: item.ProductName(), Quantity: int32(item.Quantity())}
	}
	return &orderspb.OrderResponse{
		OrderId: o.ID(), Status: o.Status().String(), FinalPrice: o.FinalPrice().InexactFloat64(), OrderNumber: o.OrderNumber(),
		Address: &orderspb.Address{City: addr.City, Street: addr.Street, House: addr.House, Apartment: addr.Apartment, Floor: addr.Floor, Entrance: addr.Entrance, Comment: addr.Comment},
		Items:   pbItems, ChefId: o.ChefID(), CourierId: o.CourierID(),
	}
}
