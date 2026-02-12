package handlers

import (
	"context"
	"time"
	"github.com/gofiber/fiber/v2"
	catalog_pb "github.com/versoit/diploma/services/catalog/api/proto/pb"
	"go.uber.org/zap"
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
		return HandleGrpcError(c, h.log, err, "failed to fetch products")
	}

	return SuccessResponse(c, resp.Products)
}

func (h *CatalogHandler) GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.client.GetProduct(ctx, &catalog_pb.GetProductRequest{Id: id})
	if err != nil {
		return HandleGrpcError(c, h.log, err, "product not found")
	}

	return SuccessResponse(c, resp)
}
