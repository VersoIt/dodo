package handlers

import (
	"context"
	"time"
	"github.com/gofiber/fiber/v2"
	orders_pb "github.com/versoit/diploma/services/orders/api/proto/pb"
	"go.uber.org/zap"
)

type OrderHandler struct {
	log    *zap.Logger
	client orders_pb.OrderServiceClient
}

func NewOrderHandler(log *zap.Logger, client orders_pb.OrderServiceClient) *OrderHandler {
	return &OrderHandler{
		log:    log,
		client: client,
	}
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	type OrderItem struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}

	type Request struct {
		CustomerID string      `json:"customer_id"`
		City       string      `json:"city"`
		Street     string      `json:"street"`
		Items      []OrderItem `json:"items"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "invalid request")
	}

	h.log.Info("Creating new order via gRPC")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pbItems := make([]*orders_pb.OrderItem, len(req.Items))
	for i, item := range req.Items {
		pbItems[i] = &orders_pb.OrderItem{
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
		}
	}

	resp, err := h.client.CreateOrder(ctx, &orders_pb.CreateOrderRequest{
		CustomerId: req.CustomerID,
		Address: &orders_pb.Address{
			City:   req.City,
			Street: req.Street,
		},
		Items: pbItems,
	})

	if err != nil {
		return HandleGrpcError(c, h.log, err, "failed to create order")
	}

	return SuccessResponse(c, resp)
}

func (h *OrderHandler) PayOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	h.log.Info("Paying for order", zap.String("order_id", id))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.PayOrder(ctx, &orders_pb.PayOrderRequest{
		OrderId: id,
	})
	if err != nil {
		return HandleGrpcError(c, h.log, err, "failed to pay for order")
	}

	return SuccessResponse(c, resp)
}

func (h *OrderHandler) GetOrderStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	h.log.Info("Checking order status", zap.String("order_id", id))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.GetOrder(ctx, &orders_pb.GetOrderRequest{
		OrderId: id,
	})
	if err != nil {
		return HandleGrpcError(c, h.log, err, "order not found")
	}

	return SuccessResponse(c, resp)
}
