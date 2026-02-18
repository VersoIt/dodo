package grpc

import (
	"context"

	logistics_pb "github.com/versoit/diploma/be/logistics/api/proto/pb"
	"github.com/versoit/diploma/be/orders/internal/config"
	"github.com/versoit/diploma/be/orders/internal/domain"
	"go.uber.org/fx"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type logisticsClient struct {
	client logistics_pb.DeliveryServiceClient
}

func NewLogisticsClient(lc fx.Lifecycle, cfg *config.Config) (domain.LogisticsService, error) {
	addr := cfg.Services.Logistics

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
