package app

import (
	"github.com/versoit/diploma/services/auth/internal/api/grpc"
	"github.com/versoit/diploma/services/auth/internal/repository"
	"github.com/versoit/diploma/services/auth/usecase"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		repository.NewInMemoryUserRepository,
		usecase.NewAuthUseCase,
		grpc.NewAuthHandler,
	),
)
