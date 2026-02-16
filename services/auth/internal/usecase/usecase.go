package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/versoit/diploma/services/auth/internal/domain"
)

var (
	ErrUnauthorized = errors.New("unauthorized: invalid email or password")
	ErrUserExists   = errors.New("user with this email already exists")
	ErrInvalidInput = errors.New("invalid input data")
)

type AuthUseCase struct {
	repo domain.UserRepository
	tm   trm.Manager
	log  *slog.Logger
}

func NewAuthUseCase(repo domain.UserRepository, tm trm.Manager, log *slog.Logger) *AuthUseCase {
	return &AuthUseCase{
		repo: repo,
		tm:   tm,
		log:  log,
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, email, password, name string, role domain.Role) (*domain.User, error) {
	if email == "" || password == "" {
		return nil, fmt.Errorf("%w: email and password are required", ErrInvalidInput)
	}

	uc.log.Info("registering user", slog.String("email", email), slog.String("name", name), slog.String("role", role.String()))

	var user *domain.User
	err := uc.tm.Do(ctx, func(ctx context.Context) error {
		existing, err := uc.repo.FindByEmail(ctx, email)
		if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
			return fmt.Errorf("failed to check existing user: %w", err)
		}
		if existing != nil {
			uc.log.Warn("registration failed: user already exists", slog.String("email", email))
			return ErrUserExists
		}

		user, err = domain.NewUser(email, password, name, role)
		if err != nil {
			return fmt.Errorf("domain validation failed: %w", err)
		}

		if err = uc.repo.Save(ctx, user); err != nil {
			return fmt.Errorf("failed to save new user: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	uc.log.Info("user registered successfully", slog.String("user_id", user.ID()))
	return user, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (*domain.User, error) {
	if email == "" || password == "" {
		return nil, fmt.Errorf("%w: credentials required", ErrInvalidInput)
	}

	uc.log.Info("login attempt", slog.String("email", email))

	user, err := uc.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
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

func (uc *AuthUseCase) GetUser(ctx context.Context, id string) (*domain.User, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: user ID is required", ErrInvalidInput)
	}

	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("retrieve user: %w", err)
	}

	return user, nil
}

func (uc *AuthUseCase) UpdateUser(ctx context.Context, id, name, phone string) (*domain.User, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: user ID is required", ErrInvalidInput)
	}

	var user *domain.User
	err := uc.tm.Do(ctx, func(ctx context.Context) error {
		var err error
		user, err = uc.repo.FindByID(ctx, id)
		if err != nil {
			return fmt.Errorf("find user: %w", err)
		}

		user.UpdateProfile(name, phone)

		if err := uc.repo.Save(ctx, user); err != nil {
			return fmt.Errorf("save updated user: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	uc.log.Info("user profile updated", slog.String("user_id", id))
	return user, nil
}
