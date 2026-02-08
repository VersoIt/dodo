package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/versoit/diploma/services/notification"
)

var (
	ErrInvalidInput = errors.New("invalid input data")
)

type NotificationUseCase struct {
	repo notification.NotificationRepository
}

func NewNotificationUseCase(repo notification.NotificationRepository) *NotificationUseCase {
	return &NotificationUseCase{repo: repo}
}

func (uc *NotificationUseCase) NotifyUser(ctx context.Context, userID string, title, msg string) error {
	if userID == "" || title == "" || msg == "" {
		return fmt.Errorf("%w: user ID, title and message are mandatory", ErrInvalidInput)
	}

	n := notification.NewNotification(userID, notification.ChannelPush, title, msg)

	// Здесь может быть вызов внешнего API. В случае ошибки мы все равно можем захотеть
	// сохранить запись о неудачном уведомлении.
	n.MarkSent()

	if err := uc.repo.Save(ctx, n); err != nil {
		return fmt.Errorf("failed to persist notification log: %w", err)
	}

	return nil
}