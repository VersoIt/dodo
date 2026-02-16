package grpc

import (
	"context"
	"os"

	"github.com/versoit/diploma/pkg/common"
	treasury_pb "github.com/versoit/diploma/services/treasury/api/proto/pb"
	"github.com/versoit/diploma/services/orders/internal/domain"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type treasuryClient struct {
	client treasury_pb.PaymentServiceClient
}

func NewTreasuryClient(lc fx.Lifecycle) (domain.TreasuryService, error) {
	addr := os.Getenv("TREASURY_SERVICE_ADDR")
	if addr == "" {
		addr = "treasury:8080"
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
