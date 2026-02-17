package grpc

import (
	"context"
	"os"

	logistics_pb "github.com/versoit/diploma/be/logistics/api/proto/pb"
	"github.com/versoit/diploma/be/orders/internal/domain"
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

func (c *logisticsClient) CreateDelivery(ctx context.Context, orderID string, orderNumber string, address domain.DeliveryAddress, items []*domain.OrderItem) error {
	reqItems := make([]*logistics_pb.DeliveryItem, len(items))
	for i, item := range items {
		reqItems[i] = &logistics_pb.DeliveryItem{
			ProductId:   item.ProductID(),
			ProductName: item.ProductName(),
			Quantity:    int32(item.Quantity()),
		}
	}

	_, err := c.client.CreateDelivery(ctx, &logistics_pb.CreateDeliveryRequest{
		OrderId:     orderID,
		OrderNumber: orderNumber,
		City:        address.City,
		Street:      address.Street,
		House:       address.House,
		Apartment:   address.Apartment,
		Items:       reqItems,
	})
	return err
}
