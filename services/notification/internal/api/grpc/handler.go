package grpc

import (
	"context"

	"github.com/versoit/diploma/services/notification/api/proto/pb"
	"github.com/versoit/diploma/services/notification/usecase"
	"google.golang.org/grpc"
)

type NotificationHandler struct {
	pb.UnimplementedNotificationServiceServer
	uc *usecase.NotificationUseCase
}

func NewNotificationHandler(uc *usecase.NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{uc: uc}
}

func (h *NotificationHandler) Register(server *grpc.Server) {
	pb.RegisterNotificationServiceServer(server, h)
}

func (h *NotificationHandler) SendNotification(ctx context.Context, req *pb.NotificationRequest) (*pb.NotificationResponse, error) {
	err := h.uc.NotifyUser(ctx, req.UserId, req.Title, req.Message)
	if err != nil {
		return nil, err
	}

	return &pb.NotificationResponse{Success: true}, nil
}
