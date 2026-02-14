package grpc

import (
	"context"
	"log/slog"

	notification_pb "github.com/versoit/diploma/services/notification/api/proto/pb"
	"github.com/versoit/diploma/services/notification/usecase"
	"google.golang.org/grpc"
)

type NotificationHandler struct {
	notification_pb.UnimplementedNotificationServiceServer
	uc  *usecase.NotificationUseCase
	log *slog.Logger
}

func NewNotificationHandler(uc *usecase.NotificationUseCase, log *slog.Logger) *NotificationHandler {
	return &NotificationHandler{
		uc:  uc,
		log: log,
	}
}

func (h *NotificationHandler) Register(server *grpc.Server) {
	notification_pb.RegisterNotificationServiceServer(server, h)
}

func (h *NotificationHandler) SendNotification(ctx context.Context, req *notification_pb.NotificationRequest) (*notification_pb.NotificationResponse, error) {
	h.log.Info("Sending notification", slog.String("user_id", req.UserId), slog.String("title", req.Title))
	
	err := h.uc.NotifyUser(ctx, req.UserId, req.Title, req.Message)
	if err != nil {
		h.log.Error("Failed to send notification", slog.String("user_id", req.UserId), slog.Any("error", err))
		return nil, err
	}

	return &notification_pb.NotificationResponse{Success: true}, nil
}
