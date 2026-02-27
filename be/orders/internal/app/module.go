package app

import (
	"context"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/versoit/diploma/be/orders/internal/config"
	"github.com/versoit/diploma/be/orders/internal/handler"
	infrastructure "github.com/versoit/diploma/be/orders/internal/infrastructure/grpc"
	"github.com/versoit/diploma/be/orders/internal/infrastructure/outbox"
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
		func(cfg *config.Config) (*nats.Conn, error) {
			return nats.Connect(cfg.Nats.URL)
		},
		common.NewPGXPool,
		telemetry.NewTracerProvider,
		common.NewTransactionManager,
		repository.NewOrderRepository,
		outbox.NewRelay,
		infrastructure.NewCatalogClient,
		infrastructure.NewKitchenClient,
		infrastructure.NewLogisticsClient,
		infrastructure.NewTreasuryClient,
		infrastructure.NewNotificationClient,
		usecase.NewOrderUseCase,
		handler.NewOrdersHandler,
	),
	fx.Invoke(
		func(*sdktrace.TracerProvider) {},
		func(lc fx.Lifecycle, r *outbox.Relay, log *slog.Logger) {
			ctx, cancel := context.WithCancel(context.Background())
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					log.Info("Starting Outbox Relay")
					go r.Start(ctx)
					return nil
				},
				OnStop: func(_ context.Context) error {
					log.Info("Stopping Outbox Relay")
					cancel()
					return nil
				},
			})
		},
	),
)
