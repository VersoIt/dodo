package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/versoit/diploma/services/auth"
)

var (
	ErrUnauthorized = errors.New("unauthorized: invalid email or password")
)

type AuthUseCase struct {
	repo auth.UserRepository
}

func NewAuthUseCase(repo auth.UserRepository) *AuthUseCase {
	return &AuthUseCase{repo: repo}
}

// Register создает нового пользователя.
func (uc *AuthUseCase) Register(ctx context.Context, email, password string, role auth.Role) (*auth.User, error) {
	// Проверка на существование
	existing, _ := uc.repo.FindByEmail(email)
	if existing != nil {
		return nil, errors.New("user already exists")
	}

	// Создание через доменную фабрику
	user, err := auth.NewUser(email, password, role)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return user, nil
}

// Login проверяет данные и возвращает пользователя (без генерации JWT здесь, это задача API слоя).
func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (*auth.User, error) {
	user, err := uc.repo.FindByEmail(email)
	if err != nil {
		return nil, ErrUnauthorized
	}

	if !user.CheckPassword(password) {
		return nil, ErrUnauthorized
	}

	return user, nil
}
