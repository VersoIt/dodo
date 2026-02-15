package handlers

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	orders_pb "github.com/versoit/diploma/services/orders/api/proto/pb"
)

type OrderHandler struct {
	log    *slog.Logger
	client orders_pb.OrderServiceClient
}

func NewOrderHandler(log *slog.Logger, client orders_pb.OrderServiceClient) *OrderHandler {
	return &OrderHandler{log: log, client: client}
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return ErrorResponse(c, fiber.StatusUnauthorized, "unauthorized")
	}

	var req struct {
		Items []struct {
			ProductID string `json:"product_id"`
			Quantity  int    `json:"quantity"`
		} `json:"items"`
		Address struct {
			City      string `json:"city"`
			Street    string `json:"street"`
			House     string `json:"house"`
			Apartment string `json:"apartment"`
			Floor     string `json:"floor"`
			Entrance  string `json:"entrance"`
			Comment   string `json:"comment"`
		} `json:"address"`
		PromoCode string `json:"promo_code"`
	}

	if err := c.BodyParser(&req); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	pbItems := make([]*orders_pb.OrderItem, len(req.Items))
	for i, it := range req.Items {
		pbItems[i] = &orders_pb.OrderItem{
			ProductId: it.ProductID,
			Quantity:  int32(it.Quantity),
		}
	}

	ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.client.CreateOrder(ctx, &orders_pb.CreateOrderRequest{
		CustomerId: userID,
		Items:      pbItems,
		Address: &orders_pb.Address{
			City:      req.Address.City,
			Street:    req.Address.Street,
			House:     req.Address.House,
			Apartment: req.Address.Apartment,
			Floor:     req.Address.Floor,
			Entrance:  req.Address.Entrance,
			Comment:   req.Address.Comment,
		},
		PromoCode: req.PromoCode,
	})

	if err != nil {
		return HandleGrpcError(c, h.log, err, "failed to create order")
	}

	return SuccessResponse(c, resp)
}

func (h *OrderHandler) GetOrderStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := h.client.GetOrder(ctx, &orders_pb.GetOrderRequest{OrderId: id})
	if err != nil { return HandleGrpcError(c, h.log, err, "not found") }
	return SuccessResponse(c, resp)
}

func (h *OrderHandler) ListOrders(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := h.client.ListOrders(ctx, &orders_pb.ListOrdersRequest{CustomerId: userID})
	if err != nil { return HandleGrpcError(c, h.log, err, "failed to list") }
	return SuccessResponse(c, resp.Orders)
}

func (h *OrderHandler) ListAllOrders(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := h.client.ListAllOrders(ctx, &orders_pb.ListAllOrdersRequest{})
	if err != nil { return HandleGrpcError(c, h.log, err, "failed to list") }
	return SuccessResponse(c, resp.Orders)
}

func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	id := c.Params("id"); userID := c.Locals("user_id").(string)
	var req struct { Status string `json:"status"` }
	if err := c.BodyParser(&req); err != nil { return ErrorResponse(c, fiber.StatusBadRequest, "invalid request") }
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := h.client.UpdateOrderStatus(ctx, &orders_pb.UpdateOrderStatusRequest{OrderId: id, Status: req.Status, PerformerId: userID})
	if err != nil { return HandleGrpcError(c, h.log, err, "failed update") }
	return SuccessResponse(c, resp)
}

func (h *OrderHandler) PayOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := h.client.PayOrder(ctx, &orders_pb.PayOrderRequest{OrderId: id})
	if err != nil { return HandleGrpcError(c, h.log, err, "payment failed") }
	return SuccessResponse(c, resp)
}

// --- Promo Codes ---

func (h *OrderHandler) CreatePromoCode(c *fiber.Ctx) error {
	var req struct { Code string `json:"code"`; Type string `json:"type"`; Amount float64 `json:"amount"` }
	if err := c.BodyParser(&req); err != nil { return ErrorResponse(c, fiber.StatusBadRequest, "invalid request") }
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := h.client.CreatePromoCode(ctx, &orders_pb.CreatePromoCodeRequest{Code: req.Code, Type: req.Type, Amount: req.Amount})
	if err != nil { return HandleGrpcError(c, h.log, err, "failed to create promo") }
	return SuccessResponse(c, resp)
}

func (h *OrderHandler) ListPromos(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := h.client.ListPromos(ctx, &orders_pb.ListPromosRequest{})
	if err != nil { return HandleGrpcError(c, h.log, err, "failed to list promos") }
	return SuccessResponse(c, resp.Promos)
}

func (h *OrderHandler) DeletePromo(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := h.client.DeletePromo(ctx, &orders_pb.DeletePromoRequest{Id: id})
	if err != nil { return HandleGrpcError(c, h.log, err, "failed to delete") }
	return SuccessResponse(c, resp)
}

func (h *OrderHandler) CheckPromoCode(c *fiber.Ctx) error {
	code := c.Params("code")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := h.client.CheckPromoCode(ctx, &orders_pb.CheckPromoCodeRequest{Code: code})
	if err != nil { return HandleGrpcError(c, h.log, err, "promo not found") }
	return SuccessResponse(c, resp)
}

// --- Analytics ---

func (h *OrderHandler) GetAnalytics(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := h.client.GetAnalytics(ctx, &orders_pb.GetAnalyticsRequest{})
	if err != nil { return HandleGrpcError(c, h.log, err, "failed to get analytics") }
	return SuccessResponse(c, resp)
}
