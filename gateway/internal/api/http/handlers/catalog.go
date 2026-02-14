package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	catalog_pb "github.com/versoit/diploma/services/catalog/api/proto/pb"
	"log/slog"
)

type CatalogHandler struct {
	log    *slog.Logger
	client catalog_pb.ProductServiceClient
}

func NewCatalogHandler(log *slog.Logger, client catalog_pb.ProductServiceClient) *CatalogHandler {
	return &CatalogHandler{
		log:    log,
		client: client,
	}
}

func (h *CatalogHandler) ListProducts(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.ListProducts(ctx, &catalog_pb.ListProductsRequest{})
	if err != nil {
		return HandleGrpcError(c, h.log, err, "failed to list products")
	}

	return SuccessResponse(c, resp)
}

func (h *CatalogHandler) GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.client.GetProduct(ctx, &catalog_pb.GetProductRequest{
		Id: id,
	})
	if err != nil {
		return HandleGrpcError(c, h.log, err, "product not found")
	}

	return SuccessResponse(c, resp)
}
