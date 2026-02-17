package app

import (
	"context"

	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/be/treasury/internal/handler"
	"github.com/versoit/diploma/be/treasury/internal/repository"
	"github.com/versoit/diploma/be/treasury/internal/usecase"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		func() (context.Context, context.CancelFunc) {
			return context.WithCancel(context.Background())
		},
		common.NewPGXPool,
		common.NewTransactionManager,
		repository.NewPaymentRepository,
		usecase.NewTreasuryUseCase,
		handler.NewTreasuryHandler,
	),
)
