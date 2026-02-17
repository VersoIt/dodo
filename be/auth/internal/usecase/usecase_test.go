package usecase

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/versoit/diploma/be/auth/internal/domain"
)

type MockUserRepo struct {
	usersByEmail map[string]*domain.User
	usersByID    map[string]*domain.User
}

func NewMockUserRepo() *MockUserRepo {
	return &MockUserRepo{
		usersByEmail: make(map[string]*domain.User),
		usersByID:    make(map[string]*domain.User),
	}
}

func (m *MockUserRepo) Save(ctx context.Context, u *domain.User) error {
	m.usersByEmail[u.Email()] = u
	m.usersByID[u.ID()] = u
	return nil
}

func (m *MockUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	if u, ok := m.usersByEmail[email]; ok {
		return u, nil
	}
	return nil, domain.ErrUserNotFound
}

func (m *MockUserRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	if u, ok := m.usersByID[id]; ok {
		return u, nil
	}
	return nil, domain.ErrUserNotFound
}

type dummyTM struct{}

func (d *dummyTM) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func (d *dummyTM) DoWithSettings(ctx context.Context, _ trm.Settings, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestAuthUseCase_Register(t *testing.T) {
	repo := NewMockUserRepo()
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	uc := NewAuthUseCase(repo, &dummyTM{}, log)

	user, err := uc.Register(context.Background(), "test@example.com", "password123", "Test User", domain.RoleClient)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Email() != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", user.Email())
	}

	_, err = uc.Register(context.Background(), "test@example.com", "password456", "Test User 2", domain.RoleClient)
	if err != ErrUserExists {
		t.Errorf("expected ErrUserExists, got %v", err)
	}
}

func TestAuthUseCase_Login(t *testing.T) {
	repo := NewMockUserRepo()
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	uc := NewAuthUseCase(repo, &dummyTM{}, log)

	if _, err := uc.Register(context.Background(), "user@example.com", "secret123", "User", domain.RoleClient); err != nil {
		t.Fatalf("setup failed: failed to register user: %v", err)
	}

	user, err := uc.Login(context.Background(), "user@example.com", "secret123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user == nil {
		t.Error("expected user, got nil")
	}

	_, err = uc.Login(context.Background(), "user@example.com", "wrongpass")
	if err != ErrUnauthorized {
		t.Errorf("expected unauthorized error, got %v", err)
	}
}
