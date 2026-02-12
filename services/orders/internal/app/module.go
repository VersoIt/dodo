package app

import (
	"context"

	"github.com/versoit/diploma/pkg/common"
		"github.com/versoit/diploma/services/orders/internal/api/grpc"
		infrastructure "github.com/versoit/diploma/services/orders/internal/infrastructure/grpc"
		"github.com/versoit/diploma/services/orders/internal/repository"
		"github.com/versoit/diploma/services/orders/usecase"
		"go.uber.org/fx"
	)
	
	var Module = fx.Options(
		fx.Provide(
			func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
					common.NewPGXPool,
					repository.NewOrderRepository,
					infrastructure.NewCatalogClient,
					infrastructure.NewAnalyticsClient,
					usecase.NewOrderUseCase,
			
			grpc.NewOrdersHandler,
		),
	)
	