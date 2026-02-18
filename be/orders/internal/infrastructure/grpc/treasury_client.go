package grpc

import (
	"context"

	"github.com/versoit/diploma/be/orders/internal/config"
	"github.com/versoit/diploma/be/orders/internal/domain"
	"github.com/versoit/diploma/pkg/common"
	treasury_pb "github.com/versoit/diploma/be/treasury/api/proto/pb"
	"go.uber.org/fx"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type treasuryClient struct {
	client treasury_pb.PaymentServiceClient
}

func NewTreasuryClient(lc fx.Lifecycle, cfg *config.Config) (domain.TreasuryService, error) {
	addr := cfg.Services.Treasury

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

	return &treasuryClient{client: treasury_pb.NewPaymentServiceClient(conn)}, nil
}

func (c *treasuryClient) ProcessPayment(ctx context.Context, orderID string, amount common.Money) error {
	_, err := c.client.ProcessPayment(ctx, &treasury_pb.PaymentRequest{
		OrderId: orderID,
		Amount:  amount.InexactFloat64(),
		Method:  1, // Default method
	})
	return err
}
