package repository

import (
	"context"
	"sync"

	"github.com/versoit/diploma/services/notification"
)

type InMemoryNotificationRepository struct {
	mu    sync.RWMutex
	store []*notification.Notification
}

func NewInMemoryNotificationRepository() notification.NotificationRepository {
	return &InMemoryNotificationRepository{
		store: make([]*notification.Notification, 0),
	}
}

func (r *InMemoryNotificationRepository) Save(ctx context.Context, n *notification.Notification) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store = append(r.store, n)
	return nil
}
