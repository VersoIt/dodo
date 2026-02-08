package app

import (
	"context"

	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/services/treasury/internal/api/grpc"
	"github.com/versoit/diploma/services/treasury/internal/repository"
	"github.com/versoit/diploma/services/treasury/usecase"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		func() (context.Context, context.CancelFunc) {
			return context.WithCancel(context.Background())
		},
		common.NewPGXPool,
		repository.NewPaymentRepository,
		usecase.NewTreasuryUseCase,
		grpc.NewTreasuryHandler,
	),
)