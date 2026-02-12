package grpc

import (
	"context"
	"os"

	"github.com/versoit/diploma/services/orders"
	analytics_pb "github.com/versoit/diploma/services/analytics/api/proto/pb"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type analyticsClient struct {
	client analytics_pb.KpiServiceClient
}

func NewAnalyticsClient(lc fx.Lifecycle) (orders.AnalyticsService, error) {
	addr := os.Getenv("ANALYTICS_SERVICE_ADDR")
	if addr == "" {
		addr = "analytics:8080"
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

	return &analyticsClient{client: analytics_pb.NewKpiServiceClient(conn)}, nil
}

func (c *analyticsClient) ReportSale(ctx context.Context, managerID string, amount float64) error {
	_, err := c.client.RecordSale(ctx, &analytics_pb.RecordSaleRequest{
		ManagerId: managerID,
		Amount:    amount,
	})
	return err
}
