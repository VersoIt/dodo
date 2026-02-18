package grpc

import (
	"context"

	kitchen_pb "github.com/versoit/diploma/be/kitchen/api/proto/pb"
	"github.com/versoit/diploma/be/orders/internal/config"
	"github.com/versoit/diploma/be/orders/internal/domain"
	"go.uber.org/fx"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type kitchenClient struct {
	client kitchen_pb.TicketServiceClient
}

func NewKitchenClient(lc fx.Lifecycle, cfg *config.Config) (domain.KitchenService, error) {
	addr := cfg.Services.Kitchen

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

	return &kitchenClient{client: kitchen_pb.NewTicketServiceClient(conn)}, nil
}

func (c *kitchenClient) CreateTicket(ctx context.Context, orderID string, orderNumber string, items []*domain.OrderItem) error {
	reqItems := make([]*kitchen_pb.KitchenItem, len(items))

	for i, item := range items {
		reqItems[i] = &kitchen_pb.KitchenItem{
			ProductId:   item.ProductID(),
			ProductName: item.ProductName(),
			Quantity:    int32(item.Quantity()),
		}
	}

	_, err := c.client.CreateTicket(ctx, &kitchen_pb.CreateTicketRequest{
		OrderId:     orderID,
		OrderNumber: orderNumber,
		Items:       reqItems,
	})
	return err
}
