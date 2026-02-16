package grpc

import (
	"context"
	"fmt"
	"os"

	notification_pb "github.com/versoit/diploma/services/notification/api/proto/pb"
	"github.com/versoit/diploma/services/orders/internal/domain"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type notificationClient struct {
	client notification_pb.NotificationServiceClient
}

func NewNotificationClient(lc fx.Lifecycle) (domain.NotificationService, error) {
	addr := os.Getenv("NOTIFICATION_SERVICE_ADDR")
	if addr == "" {
		addr = "notification:8080"
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
