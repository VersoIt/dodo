package grpc

import (
	"context"

	"github.com/versoit/diploma/gateway/internal/config"
	auth_pb "github.com/versoit/diploma/be/auth/api/proto/pb"
	catalog_pb "github.com/versoit/diploma/be/catalog/api/proto/pb"
	kitchen_pb "github.com/versoit/diploma/be/kitchen/api/proto/pb"
	logistics_pb "github.com/versoit/diploma/be/logistics/api/proto/pb"
	orders_pb "github.com/versoit/diploma/be/orders/api/proto/pb"
	"go.uber.org/fx"
	"log/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Module = fx.Provide(
	NewAuthClient,
	NewCatalogClient,
	NewOrdersClient,
	NewKitchenClient,
	NewLogisticsClient,
)

func NewAuthClient(lc fx.Lifecycle, cfg *config.Config, log *slog.Logger) (auth_pb.UserServiceClient, error) {
	conn, err := dial(lc, cfg.AuthService)
	if err != nil {
		return nil, err
	}
	return auth_pb.NewUserServiceClient(conn), nil
}

func NewCatalogClient(lc fx.Lifecycle, cfg *config.Config, log *slog.Logger) (catalog_pb.ProductServiceClient, error) {
	conn, err := dial(lc, cfg.CatalogService)
	if err != nil {
		return nil, err
	}
	return catalog_pb.NewProductServiceClient(conn), nil
}

func NewOrdersClient(lc fx.Lifecycle, cfg *config.Config, log *slog.Logger) (orders_pb.OrderServiceClient, error) {
	conn, err := dial(lc, cfg.OrdersService)
	if err != nil {
		return nil, err
	}
	return orders_pb.NewOrderServiceClient(conn), nil
}

func NewKitchenClient(lc fx.Lifecycle, cfg *config.Config, log *slog.Logger) (kitchen_pb.TicketServiceClient, error) {
	conn, err := dial(lc, cfg.KitchenService)
	if err != nil {
		return nil, err
	}
	return kitchen_pb.NewTicketServiceClient(conn), nil
}

func NewLogisticsClient(lc fx.Lifecycle, cfg *config.Config, log *slog.Logger) (logistics_pb.DeliveryServiceClient, error) {
	conn, err := dial(lc, cfg.LogisticsService)
	if err != nil {
		return nil, err
	}
	return logistics_pb.NewDeliveryServiceClient(conn), nil
}

func dial(lc fx.Lifecycle, addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return conn.Close()
		},
	})
	return conn, nil
}
