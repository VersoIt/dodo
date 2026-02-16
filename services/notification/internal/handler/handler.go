package handler

import (
	"context"
	"log/slog"

	notificationpb "github.com/versoit/diploma/services/notification/api/proto/pb"
	"github.com/versoit/diploma/services/notification/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotificationHandler struct {
	notificationpb.UnimplementedNotificationServiceServer
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
	notificationpb.RegisterNotificationServiceServer(server, h)
}

func (h *NotificationHandler) SendNotification(ctx context.Context, req *notificationpb.NotificationRequest) (*notificationpb.NotificationResponse, error) {
	h.log.Info("Sending notification", slog.String("user_id", req.UserId), slog.String("title", req.Title))

	err := h.uc.NotifyUser(ctx, req.UserId, req.Title, req.Message)
	if err != nil {
		h.log.Error("Failed to send notification", slog.String("user_id", req.UserId), slog.Any("error", err))
		return nil, status.Errorf(codes.Internal, "notify user: %v", err)
	}

	return &notificationpb.NotificationResponse{Success: true}, nil
}
