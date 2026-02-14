package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/versoit/diploma/services/auth"
)

var (
	ErrUnauthorized = errors.New("unauthorized: invalid email or password")
	ErrUserExists   = errors.New("user with this email already exists")
	ErrInvalidInput = errors.New("invalid input data")
)

type AuthUseCase struct {
	repo auth.UserRepository
	log  *slog.Logger
}

func NewAuthUseCase(repo auth.UserRepository, log *slog.Logger) *AuthUseCase {
	return &AuthUseCase{
		repo: repo,
		log:  log,
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, email, password, name string, role auth.Role) (*auth.User, error) {
	if email == "" || password == "" {
		return nil, fmt.Errorf("%w: email and password are required", ErrInvalidInput)
	}

	uc.log.Info("registering user", slog.String("email", email), slog.String("name", name), slog.String("role", role.String()))

	existing, err := uc.repo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, auth.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		uc.log.Warn("registration failed: user already exists", slog.String("email", email))
		return nil, ErrUserExists
	}

	user, err := auth.NewUser(email, password, name, role)
	if err != nil {
		return nil, fmt.Errorf("domain validation failed: %w", err)
	}

	if err := uc.repo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save new user: %w", err)
	}

	uc.log.Info("user registered successfully", slog.String("user_id", user.ID()))
	return user, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (*auth.User, error) {
	if email == "" || password == "" {
		return nil, fmt.Errorf("%w: credentials required", ErrInvalidInput)
	}

	uc.log.Info("login attempt", slog.String("email", email))

	user, err := uc.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			uc.log.Warn("login failed: user not found", slog.String("email", email))
			return nil, ErrUnauthorized
		}
		return nil, fmt.Errorf("database error during login: %w", err)
	}

	if !user.CheckPassword(password) {
		uc.log.Warn("login failed: invalid password", slog.String("email", email))
		return nil, ErrUnauthorized
	}

	uc.log.Info("user logged in", slog.String("user_id", user.ID()))
	return user, nil
}

func (uc *AuthUseCase) GetUser(ctx context.Context, id string) (*auth.User, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: user ID is required", ErrInvalidInput)
	}

	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			return nil, auth.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	return user, nil
}

func (uc *AuthUseCase) UpdateUser(ctx context.Context, id, name, phone string) (*auth.User, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: user ID is required", ErrInvalidInput)
	}

	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.UpdateProfile(name, phone)

	if err := uc.repo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	uc.log.Info("user profile updated", slog.String("user_id", id))
	return user, nil
}
