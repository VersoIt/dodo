package grpc

import (
	"context"
	"os"

	logistics_pb "github.com/versoit/diploma/services/logistics/api/proto/pb"
	"github.com/versoit/diploma/services/orders/internal/domain"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type logisticsClient struct {
	client logistics_pb.DeliveryServiceClient
}

func NewLogisticsClient(lc fx.Lifecycle) (domain.LogisticsService, error) {
	addr := os.Getenv("LOGISTICS_SERVICE_ADDR")
	if addr == "" {
		addr = "logistics:8080"
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

	return &logisticsClient{client: logistics_pb.NewDeliveryServiceClient(conn)}, nil
}

func (c *logisticsClient) CreateDelivery(ctx context.Context, orderID string) error {
	_, err := c.client.CreateDelivery(ctx, &logistics_pb.CreateDeliveryRequest{
		OrderId: orderID,
	})
	return err
}
