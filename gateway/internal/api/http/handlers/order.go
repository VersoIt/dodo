package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	orders_pb "github.com/versoit/diploma/services/orders/api/proto/pb"
	"log/slog"
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
	type Request struct {
		Items []struct {
			ProductID string `json:"product_id"`
			Quantity  int32  `json:"quantity"`
		} `json:"items"`
		Address struct {
			City   string `json:"city"`
			Street string `json:"street"`
		} `json:"address"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "invalid request")
	}

	userID := c.Locals("user_id").(string)

	items := make([]*orders_pb.OrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = &orders_pb.OrderItem{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.client.CreateOrder(ctx, &orders_pb.CreateOrderRequest{
		CustomerId: userID,
		Items:      items,
		Address: &orders_pb.Address{
			City:   req.Address.City,
			Street: req.Address.Street,
		},
	})
	if err != nil {
		return HandleGrpcError(c, h.log, err, "failed to create order")
	}

	return SuccessResponse(c, resp)
}

func (h *OrderHandler) PayOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	h.log.Info("Paying for order", slog.String("order_id", id))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.client.PayOrder(ctx, &orders_pb.PayOrderRequest{
		OrderId: id,
	})
	if err != nil {
		return HandleGrpcError(c, h.log, err, "payment failed")
	}

	return SuccessResponse(c, resp)
}

func (h *OrderHandler) GetOrderStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	h.log.Info("Checking order status", slog.String("order_id", id))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.client.GetOrder(ctx, &orders_pb.GetOrderRequest{
		OrderId: id,
	})
	if err != nil {
		return HandleGrpcError(c, h.log, err, "order not found")
	}

	return SuccessResponse(c, resp)
}

func (h *OrderHandler) ListOrders(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	h.log.Info("Listing orders for user", slog.String("user_id", userID))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.client.ListOrders(ctx, &orders_pb.ListOrdersRequest{
		CustomerId: userID,
	})
	if err != nil {
		return HandleGrpcError(c, h.log, err, "failed to list orders")
	}

	return SuccessResponse(c, resp.Orders)
}

func (h *OrderHandler) ListAllOrders(c *fiber.Ctx) error {
	h.log.Info("Listing all orders (admin/internal)")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.client.ListAllOrders(ctx, &orders_pb.ListAllOrdersRequest{})
	if err != nil {
		return HandleGrpcError(c, h.log, err, "failed to list all orders")
	}

	return SuccessResponse(c, resp.Orders)
}

func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	
	type Request struct {
		Status string `json:"status"`
	}
	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "invalid request")
	}

	h.log.Info("Updating order status", slog.String("order_id", id), slog.String("status", req.Status))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.client.UpdateOrderStatus(ctx, &orders_pb.UpdateOrderStatusRequest{
		OrderId: id,
		Status:  req.Status,
	})
	if err != nil {
		return HandleGrpcError(c, h.log, err, "failed to update status")
	}

	return SuccessResponse(c, resp)
}
