package app

import (
	"context"

	"github.com/versoit/diploma/be/orders/internal/config"
	"github.com/versoit/diploma/be/orders/internal/handler"
	infrastructure "github.com/versoit/diploma/be/orders/internal/infrastructure/grpc"
	"github.com/versoit/diploma/be/orders/internal/repository"
	"github.com/versoit/diploma/be/orders/internal/usecase"
	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/pkg/telemetry"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		config.NewConfig,
		func(cfg *config.Config) string {
			return cfg.Database.URL
		},
		func() (context.Context, context.CancelFunc) {
			return context.WithCancel(context.Background())
		},
		common.NewPGXPool,
		telemetry.NewTracerProvider,
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
	fx.Invoke(func(*sdktrace.TracerProvider) {}),
)
