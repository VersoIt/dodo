package grpc

import (
	"context"
	"fmt"

	notification_pb "github.com/versoit/diploma/be/notification/api/proto/pb"
	"github.com/versoit/diploma/be/orders/internal/config"
	"github.com/versoit/diploma/be/orders/internal/domain"
	"go.uber.org/fx"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type notificationClient struct {
	client notification_pb.NotificationServiceClient
}

func NewNotificationClient(lc fx.Lifecycle, cfg *config.Config) (domain.NotificationService, error) {
	addr := cfg.Services.Notification

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

	return &notificationClient{client: notification_pb.NewNotificationServiceClient(conn)}, nil
}

func (c *notificationClient) NotifyStatusChanged(ctx context.Context, customerID string, orderID string, status domain.OrderStatus) error {
	_, err := c.client.SendNotification(ctx, &notification_pb.NotificationRequest{
		UserId:  customerID,
		Title:   "Order Status Update",
		Message: fmt.Sprintf("Your order %s is now %s", orderID, status.String()),
	})
	return err
}
