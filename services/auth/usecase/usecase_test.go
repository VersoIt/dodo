package usecase

import (
	"context"
	"testing"

	"github.com/versoit/diploma/services/auth"
)

type MockUserRepo struct {
	usersByEmail map[string]*auth.User
	usersByID    map[string]*auth.User
}

func NewMockUserRepo() *MockUserRepo {
	return &MockUserRepo{
		usersByEmail: make(map[string]*auth.User),
		usersByID:    make(map[string]*auth.User),
	}
}

func (m *MockUserRepo) Save(u *auth.User) error {
	m.usersByEmail[u.Email()] = u
	m.usersByID[u.ID()] = u
	return nil
}

func (m *MockUserRepo) FindByEmail(email string) (*auth.User, error) {
	if u, ok := m.usersByEmail[email]; ok {
		return u, nil
	}
	return nil, auth.ErrUserNotFound
}

func (m *MockUserRepo) FindByID(id string) (*auth.User, error) {
	if u, ok := m.usersByID[id]; ok {
		return u, nil
	}
	return nil, auth.ErrUserNotFound
}

func TestAuthUseCase_Register(t *testing.T) {
	repo := NewMockUserRepo()
	uc := NewAuthUseCase(repo)

	// 1. Success
	user, err := uc.Register(context.Background(), "test@example.com", "password123", auth.RoleClient)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Email() != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", user.Email())
	}

	// 2. Duplicate email
	_, err = uc.Register(context.Background(), "test@example.com", "password456", auth.RoleClient)
	if err == nil {
		t.Error("expected error on duplicate email")
	}

	// 3. Invalid data (domain logic check via usecase)
	_, err = uc.Register(context.Background(), "invalid-email", "pass", auth.RoleClient)
	if err == nil {
		t.Error("expected error on invalid email/password")
	}
}

func TestAuthUseCase_Login(t *testing.T) {
	repo := NewMockUserRepo()
	uc := NewAuthUseCase(repo)

	// Setup
	uc.Register(context.Background(), "user@example.com", "secret123", auth.RoleClient)

	// 1. Success
	user, err := uc.Login(context.Background(), "user@example.com", "secret123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user == nil {
		t.Error("expected user, got nil")
	}

	// 2. Wrong password
	_, err = uc.Login(context.Background(), "user@example.com", "wrongpass")
	if err != ErrUnauthorized {
		t.Errorf("expected unauthorized error, got %v", err)
	}

	// 3. User not found
	_, err = uc.Login(context.Background(), "unknown@example.com", "secret123")
	if err != ErrUnauthorized {
		t.Errorf("expected unauthorized error (masked), got %v", err)
	}
}
