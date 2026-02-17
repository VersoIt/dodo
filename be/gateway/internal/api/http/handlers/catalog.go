package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	catalog_pb "github.com/versoit/diploma/be/catalog/api/proto/pb"
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.client.ListProducts(ctx, &catalog_pb.ListProductsRequest{})
	if err != nil {
		return HandleGrpcError(c, h.log, err, "failed to list products")
	}

	return SuccessResponse(c, resp.Products)
}

func (h *CatalogHandler) GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.client.GetProduct(ctx, &catalog_pb.GetProductRequest{
		Id: id,
	})
	if err != nil {
		return HandleGrpcError(c, h.log, err, "product not found")
	}

	return SuccessResponse(c, resp)
}

func (h *CatalogHandler) CreateProduct(c *fiber.Ctx) error {
	type Request struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		CategoryID  int32   `json:"category_id"`
		Price       float64 `json:"price"`
		ImageURL    string  `json:"image_url"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.client.CreateProduct(ctx, &catalog_pb.CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		CategoryId:  req.CategoryID,
		Price:       req.Price,
		ImageUrl:    req.ImageURL,
	})
	if err != nil {
		return HandleGrpcError(c, h.log, err, "failed to create product")
	}

	return SuccessResponse(c, resp)
}

func (h *CatalogHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	type Request struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		CategoryID  int32   `json:"category_id"`
		Price       float64 `json:"price"`
		ImageURL    string  `json:"image_url"`
		IsAvailable bool    `json:"is_available"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.client.UpdateProduct(ctx, &catalog_pb.UpdateProductRequest{
		Id:          id,
		Name:        req.Name,
		Description: req.Description,
		CategoryId:  req.CategoryID,
		Price:       req.Price,
		ImageUrl:    req.ImageURL,
		IsAvailable: req.IsAvailable,
	})
	if err != nil {
		return HandleGrpcError(c, h.log, err, "failed to update product")
	}

	return SuccessResponse(c, resp)
}
