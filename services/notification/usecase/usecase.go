package usecase

import (
	"context"
	"fmt"

	"github.com/versoit/diploma/services/notification"
)

type NotificationUseCase struct {
	repo notification.NotificationRepository
	// Здесь мог бы быть интерфейс отправителя (email/sms provider)
}

func NewNotificationUseCase(repo notification.NotificationRepository) *NotificationUseCase {
	return &NotificationUseCase{repo: repo}
}

// NotifyUser создает уведомление и сохраняет его (имитация отправки).
func (uc *NotificationUseCase) NotifyUser(ctx context.Context, userID string, title, msg string) error {
	// Для простоты используем Push как основной канал
	n := notification.NewNotification(userID, notification.ChannelPush, title, msg)

	// Имитируем отправку
	// В реальности здесь был бы вызов провайдера
	n.MarkSent()

	if err := uc.repo.Save(n); err != nil {
		return fmt.Errorf("failed to log notification: %w", err)
	}

	return nil
}
