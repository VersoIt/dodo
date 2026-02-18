package grpc

import (
	"context"

	"github.com/versoit/diploma/be/orders/internal/config"
	"github.com/versoit/diploma/be/orders/internal/domain"
	"github.com/versoit/diploma/pkg/common"
	catalog_pb "github.com/versoit/diploma/be/catalog/api/proto/pb"
	"go.uber.org/fx"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type catalogClient struct {
	client catalog_pb.ProductServiceClient
}

func NewCatalogClient(lc fx.Lifecycle, cfg *config.Config) (domain.CatalogService, error) {
	addr := cfg.Services.Catalog

	conn, err := grpc.NewClient(
		addr, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
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
