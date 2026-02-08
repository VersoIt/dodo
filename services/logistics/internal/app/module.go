package app

import (
	"context"

	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/services/logistics/internal/api/grpc"
	"github.com/versoit/diploma/services/logistics/internal/repository"
	"github.com/versoit/diploma/services/logistics/usecase"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		func() (context.Context, context.CancelFunc) {
			return context.WithCancel(context.Background())
		},
		common.NewPGXPool,
		repository.NewDeliveryRepository,
		repository.NewCourierRepository,
		usecase.NewLogisticsUseCase,
		grpc.NewLogisticsHandler,
	),
)