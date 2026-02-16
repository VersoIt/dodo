package handler

import (
	"context"
	"errors"
	"log/slog"

	"github.com/versoit/diploma/pkg/common"
	catalogpb "github.com/versoit/diploma/services/catalog/api/proto/pb"
	"github.com/versoit/diploma/services/catalog/internal/domain"
	"github.com/versoit/diploma/services/catalog/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CatalogHandler struct {
	catalogpb.UnimplementedProductServiceServer
	uc  *usecase.CatalogUseCase
	log *slog.Logger
}

func NewCatalogHandler(uc *usecase.CatalogUseCase, log *slog.Logger) *CatalogHandler {
	return &CatalogHandler{
		uc:  uc,
		log: log,
	}
}

func (h *CatalogHandler) Register(server *grpc.Server) {
	catalogpb.RegisterProductServiceServer(server, h)
}

func (h *CatalogHandler) CreateProduct(ctx context.Context, req *catalogpb.CreateProductRequest) (*catalogpb.ProductResponse, error) {
	h.log.Info("Creating product", slog.String("name", req.Name))

	p, err := h.uc.CreateProduct(ctx, req.Name, req.Description, domain.CategoryType(req.CategoryId), req.Price, req.ImageUrl)
	if err != nil {
		h.log.Error("Failed to create product", slog.String("name", req.Name), slog.Any("error", err))
		return nil, status.Errorf(codes.Internal, "create product: %v", err)
	}

	return h.mapProduct(p), nil
}

func (h *CatalogHandler) UpdateProduct(ctx context.Context, req *catalogpb.UpdateProductRequest) (*catalogpb.ProductResponse, error) {
	h.log.Info("Updating product", slog.String("id", req.Id))

	p, err := h.uc.UpdateProduct(
		ctx,
		req.Id,
		req.Name,
		req.Description,
		domain.CategoryType(req.CategoryId),
		req.Price,
		req.ImageUrl,
		req.IsAvailable,
	)
	if err != nil {
		h.log.Error("Failed to update product", slog.String("id", req.Id), slog.Any("error", err))
		return nil, status.Errorf(codes.Internal, "update product: %v", err)
	}

	return h.mapProduct(p), nil
}

func (h *CatalogHandler) GetProduct(ctx context.Context, req *catalogpb.GetProductRequest) (*catalogpb.ProductResponse, error) {
	p, err := h.uc.GetProduct(ctx, req.Id)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return nil, status.Errorf(codes.NotFound, "product not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "get product: %v", err)
	}

	return h.mapProduct(p), nil
}

func (h *CatalogHandler) ListProducts(ctx context.Context, req *catalogpb.ListProductsRequest) (*catalogpb.ListProductsResponse, error) {
	products, err := h.uc.ListProducts(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list products: %v", err)
	}

	pbProducts := make([]*catalogpb.ProductResponse, len(products))
	for i, p := range products {
		pbProducts[i] = h.mapProduct(p)
	}

	return &catalogpb.ListProductsResponse{
		Products: pbProducts,
	}, nil
}

func (h *CatalogHandler) mapProduct(p *domain.Product) *catalogpb.ProductResponse {
	return &catalogpb.ProductResponse{
		Id:          p.ID(),
		Name:        p.Name(),
		Description: p.Description(),
		Price:       p.BasePrice().InexactFloat64(),
		ImageUrl:    p.ImageURL(),
		CategoryId:  int32(p.Category()),
		IsAvailable: p.IsAvailable(),
	}
}
