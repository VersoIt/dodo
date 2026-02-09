package handlers

import (
	"context"
	"time"
	"github.com/gofiber/fiber/v2"
	catalog_pb "github.com/versoit/diploma/services/catalog/api/proto/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

type CatalogHandler struct {
	log    *zap.Logger
	client catalog_pb.ProductServiceClient
}

func NewCatalogHandler(log *zap.Logger, client catalog_pb.ProductServiceClient) *CatalogHandler {
	return &CatalogHandler{
		log:    log,
		client: client,
	}
}

func (h *CatalogHandler) ListProducts(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.ListProducts(ctx, &catalog_pb.ListProductsRequest{})
	if err != nil {
		return h.handleGrpcError(c, err, "failed to fetch products")
	}

	return SuccessResponse(c, resp.Products)
}

func (h *CatalogHandler) GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.GetProduct(ctx, &catalog_pb.GetProductRequest{Id: id})
	if err != nil {
		return h.handleGrpcError(c, err, "product not found")
	}

	return SuccessResponse(c, resp)
}

func (h *CatalogHandler) handleGrpcError(c *fiber.Ctx, err error, defaultMsg string) error {
	st, ok := status.FromError(err)
	if !ok {
		return ErrorResponse(c, fiber.StatusInternalServerError, defaultMsg)
	}

	h.log.Error("gRPC Error", 
		zap.String("code", st.Code().String()), 
		zap.String("msg", st.Message()),
		zap.String("request_id", c.Get("X-Request-ID")))

	switch st.Code() {
	case codes.NotFound:
		return ErrorResponse(c, fiber.StatusNotFound, st.Message())
	case codes.InvalidArgument:
		return ErrorResponse(c, fiber.StatusBadRequest, st.Message())
	case codes.Unavailable:
		return ErrorResponse(c, fiber.StatusServiceUnavailable, "service temporarily unavailable")
	default:
		return ErrorResponse(c, fiber.StatusInternalServerError, defaultMsg)
	}
}
