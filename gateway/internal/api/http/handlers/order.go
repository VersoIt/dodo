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
	return &OrderHandler{
		log:    log,
		client: client,
	}
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	
	var req struct {
		Items []struct {
			ProductID string `json:"product_id"`
			Quantity  int    `json:"quantity"`
		} `json:"items"`
		Address struct {
			Street string `json:"street"`
			City   string `json:"city"`
		} `json:"address"`
	}

	if err := c.BodyParser(&req); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	pbItems := make([]*orders_pb.OrderItem, len(req.Items))
	for i, item := range req.Items {
		pbItems[i] = &orders_pb.OrderItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.client.CreateOrder(ctx, &orders_pb.CreateOrderRequest{
		CustomerId: userID,
		Items:      pbItems,
		Address: &orders_pb.Address{
			Street: req.Address.Street,
			City:   req.Address.City,
		},
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
	if err != nil {
		return HandleGrpcError(c, h.log, err, "failed to get order")
	}

	return SuccessResponse(c, resp)
}

func (h *OrderHandler) ListOrders(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.client.ListOrders(ctx, &orders_pb.ListOrdersRequest{CustomerId: userID})
	if err != nil {
		return HandleGrpcError(c, h.log, err, "failed to list orders")
	}

	return SuccessResponse(c, resp.Orders)
}

func (h *OrderHandler) ListAllOrders(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.client.ListAllOrders(ctx, &orders_pb.ListAllOrdersRequest{})
	if err != nil {
		return HandleGrpcError(c, h.log, err, "failed to list orders")
	}

	return SuccessResponse(c, resp.Orders)
}

func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	
	// GET USER ID FROM LOCALS
	uID := c.Locals("user_id")
	userID, ok := uID.(string)
	if !ok || userID == "" {
		h.log.Error("CRITICAL: user_id not found in locals in UpdateOrderStatus")
		return ErrorResponse(c, fiber.StatusUnauthorized, "user identification failed")
	}
	
	type Request struct {
		Status string `json:"status"`
	}
	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "invalid request")
	}

	h.log.Info("GATEWAY: Updating status", 
		slog.String("order_id", id), 
		slog.String("status", req.Status), 
		slog.String("performer_id", userID))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.client.UpdateOrderStatus(ctx, &orders_pb.UpdateOrderStatusRequest{
		OrderId:     id,
		Status:      req.Status,
		PerformerId: userID,
	})
	if err != nil {
		return HandleGrpcError(c, h.log, err, "failed to update status")
	}

	return SuccessResponse(c, resp)
}

func (h *OrderHandler) PayOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.client.PayOrder(ctx, &orders_pb.PayOrderRequest{OrderId: id})
	if err != nil {
		return HandleGrpcError(c, h.log, err, "payment failed")
	}

	return SuccessResponse(c, resp)
}
