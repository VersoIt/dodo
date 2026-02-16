package app

import (
	"context"

	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/services/auth/internal/handler"
	"github.com/versoit/diploma/services/auth/internal/repository"
	"github.com/versoit/diploma/services/auth/internal/usecase"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		func() (context.Context, context.CancelFunc) {
			return context.WithCancel(context.Background())
		},
		common.NewPGXPool,
		common.NewTransactionManager,
		repository.NewUserRepository,
		usecase.NewAuthUseCase,
		handler.NewAuthHandler,
	),
)
