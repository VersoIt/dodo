package grpc

import (
	"context"

	notification_pb "github.com/versoit/diploma/services/notification/api/proto/pb"
	"github.com/versoit/diploma/services/notification/usecase"
	"google.golang.org/grpc"
)

type NotificationHandler struct {
	notification_pb.UnimplementedNotificationServiceServer
	uc *usecase.NotificationUseCase
}

func NewNotificationHandler(uc *usecase.NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{uc: uc}
}

func (h *NotificationHandler) Register(server *grpc.Server) {
	notification_pb.RegisterNotificationServiceServer(server, h)
}

func (h *NotificationHandler) SendNotification(ctx context.Context, req *notification_pb.NotificationRequest) (*notification_pb.NotificationResponse, error) {
	err := h.uc.NotifyUser(ctx, req.UserId, req.Title, req.Message)
	if err != nil {
		return nil, err
	}

	return &notification_pb.NotificationResponse{Success: true}, nil
}
