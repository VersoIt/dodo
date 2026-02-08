package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/versoit/diploma/services/auth"
)

var (
	ErrUnauthorized = errors.New("unauthorized: invalid email or password")
	ErrUserExists   = errors.New("user with this email already exists")
	ErrInvalidInput = errors.New("invalid input data")
)

type AuthUseCase struct {
	repo auth.UserRepository
}

func NewAuthUseCase(repo auth.UserRepository) *AuthUseCase {
	return &AuthUseCase{repo: repo}
}

func (uc *AuthUseCase) Register(ctx context.Context, email, password string, role auth.Role) (*auth.User, error) {
	if email == "" || password == "" {
		return nil, fmt.Errorf("%w: email and password are required", ErrInvalidInput)
	}

	// Проверяем, существует ли пользователь
	existing, err := uc.repo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, auth.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return nil, ErrUserExists
	}

	user, err := auth.NewUser(email, password, role)
	if err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	if err := uc.repo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save new user: %w", err)
	}

	return user, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (*auth.User, error) {
	if email == "" || password == "" {
		return nil, fmt.Errorf("%w: credentials required", ErrInvalidInput)
	}

	user, err := uc.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			return nil, ErrUnauthorized
		}
		return nil, fmt.Errorf("database error during login: %w", err)
	}

	if !user.CheckPassword(password) {
		return nil, ErrUnauthorized
	}

	return user, nil
}
