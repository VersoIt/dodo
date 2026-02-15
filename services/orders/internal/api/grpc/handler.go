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
	return &OrdersHandler{uc: uc, log: log}
}

func (h *OrdersHandler) Register(server *grpc.Server) {
	orders_pb.RegisterOrderServiceServer(server, h)
}

func (h *OrdersHandler) CreateOrder(ctx context.Context, req *orders_pb.CreateOrderRequest) (*orders_pb.OrderResponse, error) {
	h.log.Info("Creating new order", slog.String("customer_id", req.CustomerId))
	items := make([]usecase.OrderItemInput, len(req.Items))
	for i, item := range req.Items {
		items[i] = usecase.OrderItemInput{ProductID: item.ProductId, Quantity: int(item.Quantity), SizeMult: 1.0}
	}
	
	addr := orders.DeliveryAddress{
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

func (h *OrdersHandler) GetOrder(ctx context.Context, req *orders_pb.GetOrderRequest) (*orders_pb.OrderResponse, error) {
	order, err := h.uc.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found")
	}
	return h.mapOrder(order), nil
}

func (h *OrdersHandler) ListOrders(ctx context.Context, req *orders_pb.ListOrdersRequest) (*orders_pb.ListOrdersResponse, error) {
	orderList, err := h.uc.ListOrders(ctx, req.CustomerId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list orders")
	}
	pbOrders := make([]*orders_pb.OrderResponse, len(orderList))
	for i, o := range orderList {
		pbOrders[i] = h.mapOrder(o)
	}
	return &orders_pb.ListOrdersResponse{Orders: pbOrders}, nil
}

func (h *OrdersHandler) ListAllOrders(ctx context.Context, _ *orders_pb.ListAllOrdersRequest) (*orders_pb.ListOrdersResponse, error) {
	orderList, err := h.uc.ListAllOrders(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list all orders")
	}
	pbOrders := make([]*orders_pb.OrderResponse, len(orderList))
	for i, o := range orderList {
		pbOrders[i] = h.mapOrder(o)
	}
	return &orders_pb.ListOrdersResponse{Orders: pbOrders}, nil
}

func (h *OrdersHandler) UpdateOrderStatus(ctx context.Context, req *orders_pb.UpdateOrderStatusRequest) (*orders_pb.OrderResponse, error) {
	order, err := h.uc.UpdateStatus(ctx, req.OrderId, req.Status, req.PerformerId)
	if err != nil { return nil, status.Errorf(codes.Internal, "failed to update status: %v", err) }
	return h.mapOrder(order), nil
}

func (h *OrdersHandler) PayOrder(ctx context.Context, req *orders_pb.PayOrderRequest) (*orders_pb.OrderResponse, error) {
	err := h.uc.PayOrder(ctx, req.OrderId)
	if err != nil { return nil, status.Errorf(codes.Internal, "payment failed: %v", err) }
	return &orders_pb.OrderResponse{OrderId: req.OrderId, Status: "paid"}, nil
}

// --- Promo Codes ---

func (h *OrdersHandler) CreatePromoCode(ctx context.Context, req *orders_pb.CreatePromoCodeRequest) (*orders_pb.PromoCodeResponse, error) {
	p, err := h.uc.CreatePromoCode(ctx, req.Code, req.Type, req.Amount)
	if err != nil { return nil, status.Errorf(codes.Internal, "failed to create promo") }
	return &orders_pb.PromoCodeResponse{Id: p.ID(), Code: p.Code(), Type: p.DiscountType(), Amount: p.DiscountAmount().InexactFloat64(), Active: p.IsActive()}, nil
}

func (h *OrdersHandler) ListPromos(ctx context.Context, _ *orders_pb.ListPromosRequest) (*orders_pb.ListPromosResponse, error) {
	promos, _ := h.uc.ListPromos(ctx)
	pbPromos := make([]*orders_pb.PromoCodeResponse, len(promos))
	for i, p := range promos {
		pbPromos[i] = &orders_pb.PromoCodeResponse{Id: p.ID(), Code: p.Code(), Type: p.DiscountType(), Amount: p.DiscountAmount().InexactFloat64(), Active: p.IsActive()}
	}
	return &orders_pb.ListPromosResponse{Promos: pbPromos}, nil
}

func (h *OrdersHandler) DeletePromo(ctx context.Context, req *orders_pb.DeletePromoRequest) (*orders_pb.PromoDeleteResponse, error) {
	_ = h.uc.DeletePromo(ctx, req.Id)
	return &orders_pb.PromoDeleteResponse{Success: true}, nil
}

func (h *OrdersHandler) CheckPromoCode(ctx context.Context, req *orders_pb.CheckPromoCodeRequest) (*orders_pb.PromoCodeResponse, error) {
	p, err := h.uc.GetPromoByCode(ctx, req.Code)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "promo not found")
	}
	return &orders_pb.PromoCodeResponse{Id: p.ID(), Code: p.Code(), Type: p.DiscountType(), Amount: p.DiscountAmount().InexactFloat64(), Active: p.IsActive()}, nil
}

// --- Analytics ---

func (h *OrdersHandler) GetAnalytics(ctx context.Context, _ *orders_pb.GetAnalyticsRequest) (*orders_pb.AnalyticsResponse, error) {
	kpis, top, err := h.uc.GetAnalytics(ctx)
	if err != nil { return nil, err }
	
	pbTop := make([]*orders_pb.ProductStat, len(top))
	for i, s := range top {
		pbTop[i] = &orders_pb.ProductStat{Name: s.Name, Count: int32(s.Count), Revenue: s.Revenue}
	}

	return &orders_pb.AnalyticsResponse{
		TotalRevenue: kpis.TotalRevenue,
		OrdersCount:  int32(kpis.OrdersCount),
		AvgCheck:     kpis.AvgCheck,
		TopProducts:  pbTop,
	}, nil
}

func (h *OrdersHandler) mapOrder(o *orders.Order) *orders_pb.OrderResponse {
	addr := o.Address()
	pbItems := make([]*orders_pb.OrderItem, len(o.Items()))
	for i, item := range o.Items() {
		pbItems[i] = &orders_pb.OrderItem{ProductId: item.ProductID(), ProductName: item.ProductName(), Quantity: int32(item.Quantity())}
	}
	return &orders_pb.OrderResponse{
		OrderId: o.ID(), Status: o.Status().String(), FinalPrice: o.FinalPrice().InexactFloat64(), OrderNumber: o.OrderNumber(),
		Address: &orders_pb.Address{City: addr.City, Street: addr.Street, House: addr.House, Apartment: addr.Apartment, Floor: addr.Floor, Entrance: addr.Entrance, Comment: addr.Comment},
		Items: pbItems, ChefId: o.ChefID(), CourierId: o.CourierID(),
	}
}
