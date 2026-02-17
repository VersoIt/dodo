package app

import (
	"context"

	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/be/orders/internal/handler"
	infrastructure "github.com/versoit/diploma/be/orders/internal/infrastructure/grpc"
	"github.com/versoit/diploma/be/orders/internal/repository"
	"github.com/versoit/diploma/be/orders/internal/usecase"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		func() (context.Context, context.CancelFunc) {
			return context.WithCancel(context.Background())
		},
		common.NewPGXPool,
		common.NewTransactionManager,
		repository.NewOrderRepository,
		infrastructure.NewCatalogClient,
		infrastructure.NewKitchenClient,
		infrastructure.NewLogisticsClient,
		infrastructure.NewTreasuryClient,
		infrastructure.NewNotificationClient,
		usecase.NewOrderUseCase,

		handler.NewOrdersHandler,
	),
)
