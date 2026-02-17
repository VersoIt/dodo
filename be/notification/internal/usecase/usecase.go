package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/versoit/diploma/be/notification/internal/domain"
)

var (
	ErrInvalidInput = errors.New("invalid input data")
)

type NotificationUseCase struct {
	repo domain.NotificationRepository
	tm   trm.Manager
	log  *slog.Logger
}

func NewNotificationUseCase(repo domain.NotificationRepository, tm trm.Manager, log *slog.Logger) *NotificationUseCase {
	return &NotificationUseCase{
		repo: repo,
		tm:   tm,
		log:  log,
	}
}

func (uc *NotificationUseCase) NotifyUser(ctx context.Context, userID string, title, msg string) error {
	if userID == "" || title == "" || msg == "" {
		return fmt.Errorf("%w: user ID, title and message are mandatory", ErrInvalidInput)
	}

	uc.log.Info("notifying user", slog.String("user_id", userID), slog.String("title", title))

	return uc.tm.Do(ctx, func(ctx context.Context) error {
		n := domain.NewNotification(userID, domain.ChannelPush, title, msg)

		n.MarkSent()

		if err := uc.repo.Save(ctx, n); err != nil {
			uc.log.Error("failed to save notification", slog.Any("error", err), slog.String("user_id", userID))
			return fmt.Errorf("persist notification: %w", err)
		}
		return nil
	})
}
