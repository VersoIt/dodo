package grpc

import (
	"context"

	"github.com/versoit/diploma/services/catalog"
	"github.com/versoit/diploma/services/catalog/api/proto/pb"
	"github.com/versoit/diploma/services/catalog/usecase"
	"github.com/versoit/diploma/pkg/common"
	"google.golang.org/grpc"
)

type CatalogHandler struct {
	pb.UnimplementedProductServiceServer
	uc *usecase.CatalogUseCase
}

func NewCatalogHandler(uc *usecase.CatalogUseCase) *CatalogHandler {
	return &CatalogHandler{uc: uc}
}

func (h *CatalogHandler) Register(server *grpc.Server) {
	pb.RegisterProductServiceServer(server, h)
}

func (h *CatalogHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	p, err := h.uc.CreateProduct(ctx, req.Name, req.Description, catalog.CategoryType(req.CategoryId), common.Money(req.Price))
	if err != nil {
		return nil, err
	}

	return &pb.ProductResponse{
		Id:          p.ID(),
		Name:        p.Name(),
		Description: p.Description(),
		Price:       float64(p.BasePrice()),
	}, nil
}

func (h *CatalogHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	// Simple proxy
	return &pb.ProductResponse{Id: req.Id}, nil
}

func (h *CatalogHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	return &pb.ListProductsResponse{}, nil
}
