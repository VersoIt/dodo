package app

import (
	"context"

	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/be/kitchen/internal/handler"
	"github.com/versoit/diploma/be/kitchen/internal/repository"
	"github.com/versoit/diploma/be/kitchen/internal/usecase"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		func() (context.Context, context.CancelFunc) {
			return context.WithCancel(context.Background())
		},
		common.NewPGXPool,
		common.NewTransactionManager,
		repository.NewTicketRepository,
		usecase.NewKitchenUseCase,
		handler.NewKitchenHandler,
	),
)
