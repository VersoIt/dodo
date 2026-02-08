package grpc

import (
	"context"

	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/services/catalog"
	catalog_pb "github.com/versoit/diploma/services/catalog/api/proto/pb"
	"github.com/versoit/diploma/services/catalog/usecase"
	"google.golang.org/grpc"
)

type CatalogHandler struct {
	catalog_pb.UnimplementedProductServiceServer
	uc *usecase.CatalogUseCase
}

func NewCatalogHandler(uc *usecase.CatalogUseCase) *CatalogHandler {
	return &CatalogHandler{uc: uc}
}

func (h *CatalogHandler) Register(server *grpc.Server) {
	catalog_pb.RegisterProductServiceServer(server, h)
}

func (h *CatalogHandler) CreateProduct(ctx context.Context, req *catalog_pb.CreateProductRequest) (*catalog_pb.ProductResponse, error) {
	p, err := h.uc.CreateProduct(ctx, req.Name, req.Description, catalog.CategoryType(req.CategoryId), common.Money(req.Price))
	if err != nil {
		return nil, err
	}

	return &catalog_pb.ProductResponse{
		Id:          p.ID(),
		Name:        p.Name(),
		Description: p.Description(),
		Price:       float64(p.BasePrice()),
	}, nil
}

func (h *CatalogHandler) GetProduct(ctx context.Context, req *catalog_pb.GetProductRequest) (*catalog_pb.ProductResponse, error) {
	return &catalog_pb.ProductResponse{Id: req.Id}, nil
}

func (h *CatalogHandler) ListProducts(ctx context.Context, req *catalog_pb.ListProductsRequest) (*catalog_pb.ListProductsResponse, error) {
	return &catalog_pb.ListProductsResponse{}, nil
}
