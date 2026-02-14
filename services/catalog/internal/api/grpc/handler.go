package grpc

import (
	"context"
	"log/slog"

	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/services/catalog"
	catalog_pb "github.com/versoit/diploma/services/catalog/api/proto/pb"
	"github.com/versoit/diploma/services/catalog/usecase"
	"google.golang.org/grpc"
)

type CatalogHandler struct {
	catalog_pb.UnimplementedProductServiceServer
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
	catalog_pb.RegisterProductServiceServer(server, h)
}

func (h *CatalogHandler) CreateProduct(ctx context.Context, req *catalog_pb.CreateProductRequest) (*catalog_pb.ProductResponse, error) {
	h.log.Info("Creating product", slog.String("name", req.Name))
	
	p, err := h.uc.CreateProduct(ctx, req.Name, req.Description, catalog.CategoryType(req.CategoryId), common.NewMoney(req.Price))
	if err != nil {
		h.log.Error("Failed to create product", slog.String("name", req.Name), slog.Any("error", err))
		return nil, err
	}

	return &catalog_pb.ProductResponse{
		Id:          p.ID(),
		Name:        p.Name(),
		Description: p.Description(),
		Price:       p.BasePrice().InexactFloat64(),
	}, nil
}

func (h *CatalogHandler) GetProduct(ctx context.Context, req *catalog_pb.GetProductRequest) (*catalog_pb.ProductResponse, error) {
	p, err := h.uc.GetProduct(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &catalog_pb.ProductResponse{
		Id:          p.ID(),
		Name:        p.Name(),
		Description: p.Description(),
		Price:       p.BasePrice().InexactFloat64(),
	}, nil
}

func (h *CatalogHandler) ListProducts(ctx context.Context, req *catalog_pb.ListProductsRequest) (*catalog_pb.ListProductsResponse, error) {
	products, err := h.uc.ListProducts(ctx)
	if err != nil {
		return nil, err
	}

	pbProducts := make([]*catalog_pb.ProductResponse, len(products))
	for i, p := range products {
		pbProducts[i] = &catalog_pb.ProductResponse{
			Id:          p.ID(),
			Name:        p.Name(),
			Description: p.Description(),
			Price:       p.BasePrice().InexactFloat64(),
		}
	}

	return &catalog_pb.ListProductsResponse{
		Products: pbProducts,
	}, nil
}
