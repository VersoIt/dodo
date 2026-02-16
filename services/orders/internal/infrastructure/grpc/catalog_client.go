package grpc

import (
	"context"

	"github.com/versoit/diploma/services/orders/internal/domain"
	"github.com/versoit/diploma/pkg/common"
	catalog_pb "github.com/versoit/diploma/services/catalog/api/proto/pb"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

type catalogClient struct {
	client catalog_pb.ProductServiceClient
}

func NewCatalogClient(lc fx.Lifecycle) (domain.CatalogService, error) {
	addr := os.Getenv("CATALOG_SERVICE_ADDR")
	if addr == "" {
		addr = "catalog:8080"
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return conn.Close()
		},
	})

	return &catalogClient{client: catalog_pb.NewProductServiceClient(conn)}, nil
}

func (c *catalogClient) GetProduct(ctx context.Context, id string) (*domain.ProductInfo, error) {
	resp, err := c.client.GetProduct(ctx, &catalog_pb.GetProductRequest{Id: id})
	if err != nil {
		return nil, err
	}

	return &domain.ProductInfo{
		ID:        resp.Id,
		Name:      resp.Name,
		BasePrice: common.NewMoney(resp.Price),
	}, nil
}
