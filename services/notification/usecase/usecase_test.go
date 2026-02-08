package usecase

import (
	"context"
	"testing"

	"github.com/versoit/diploma/services/notification"
)

type MockNotifyRepo struct {
	lastSaved *notification.Notification
}

func (m *MockNotifyRepo) Save(n *notification.Notification) error {
	m.lastSaved = n
	return nil
}

func TestNotificationUseCase_NotifyUser(t *testing.T) {
	repo := &MockNotifyRepo{}
	uc := NewNotificationUseCase(repo)

	err := uc.NotifyUser(context.Background(), "user1", "Hello", "World")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.lastSaved == nil {
		t.Fatal("notification not saved")
	}
	if !repo.lastSaved.IsSent() {
		t.Error("notification should be marked as sent (simulated)")
	}
}
