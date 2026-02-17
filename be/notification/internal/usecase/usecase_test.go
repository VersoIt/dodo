package usecase

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/versoit/diploma/be/notification/internal/domain"
)

type MockNotifyRepo struct {
	lastSaved *domain.Notification
}

func (m *MockNotifyRepo) Save(ctx context.Context, n *domain.Notification) error {
	m.lastSaved = n
	return nil
}

type dummyTM struct{}

func (d *dummyTM) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func (d *dummyTM) DoWithSettings(ctx context.Context, _ trm.Settings, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestNotificationUseCase_NotifyUser(t *testing.T) {
	repo := &MockNotifyRepo{}
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	uc := NewNotificationUseCase(repo, &dummyTM{}, log)

	err := uc.NotifyUser(context.Background(), "user1", "Hello", "World")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.lastSaved == nil || !repo.lastSaved.IsSent() {
		t.Error("notification not marked as sent")
	}
}
