package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/versoit/diploma/services/notification"
)

var (
	ErrInvalidInput = errors.New("invalid input data")
)

type NotificationUseCase struct {
	repo notification.NotificationRepository
	log  *slog.Logger
}

func NewNotificationUseCase(repo notification.NotificationRepository, log *slog.Logger) *NotificationUseCase {
	return &NotificationUseCase{
		repo: repo,
		log:  log,
	}
}

func (uc *NotificationUseCase) NotifyUser(ctx context.Context, userID string, title, msg string) error {
	if userID == "" || title == "" || msg == "" {
		return fmt.Errorf("%w: user ID, title and message are mandatory", ErrInvalidInput)
	}

	uc.log.Info("notifying user", slog.String("user_id", userID), slog.String("title", title))

	n := notification.NewNotification(userID, notification.ChannelPush, title, msg)

	n.MarkSent()

	if err := uc.repo.Save(ctx, n); err != nil {
		uc.log.Error("failed to save notification", slog.Any("error", err), slog.String("user_id", userID))
		return fmt.Errorf("failed to persist notification log: %w", err)
	}

	return nil
}
