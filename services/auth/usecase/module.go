package usecase

import (
	"github.com/versoit/diploma/services/auth/internal/api/grpc"
	"github.com/versoit/diploma/services/auth/internal/repository"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		repository.NewInMemoryUserRepository,
		NewAuthUseCase,
		grpc.NewAuthHandler,
	),
)
