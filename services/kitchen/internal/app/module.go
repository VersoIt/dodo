package app

import (
	"context"

	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/services/kitchen/internal/api/grpc"
	"github.com/versoit/diploma/services/kitchen/internal/repository"
	"github.com/versoit/diploma/services/kitchen/usecase"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		func() (context.Context, context.CancelFunc) {
			return context.WithCancel(context.Background())
		},
		common.NewPGXPool,
		repository.NewTicketRepository,
		usecase.NewKitchenUseCase,
		grpc.NewKitchenHandler,
	),
)
