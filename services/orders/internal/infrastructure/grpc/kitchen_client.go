package grpc

import (
	"context"
	"os"

	kitchen_pb "github.com/versoit/diploma/services/kitchen/api/proto/pb"
	"github.com/versoit/diploma/services/orders/internal/domain"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type kitchenClient struct {
	client kitchen_pb.TicketServiceClient
}

func NewKitchenClient(lc fx.Lifecycle) (domain.KitchenService, error) {
	addr := os.Getenv("KITCHEN_SERVICE_ADDR")
	if addr == "" {
		addr = "kitchen:8080"
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

	return &kitchenClient{client: kitchen_pb.NewTicketServiceClient(conn)}, nil
}

func (c *kitchenClient) CreateTicket(ctx context.Context, orderID string, items []*domain.OrderItem) error {
	reqItems := make([]*kitchen_pb.KitchenItem, len(items))

	for i, item := range items {
		reqItems[i] = &kitchen_pb.KitchenItem{
			ProductId:   item.ProductID(),
			ProductName: item.ProductName(),
			Quantity:    int32(item.Quantity()),
		}
	}

	_, err := c.client.CreateTicket(ctx, &kitchen_pb.CreateTicketRequest{
		OrderId: orderID,
		Items:   reqItems,
	})
	return err
}
