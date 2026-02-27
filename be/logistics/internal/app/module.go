package app

import (
	"context"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/versoit/diploma/be/logistics/internal/config"
	"github.com/versoit/diploma/be/logistics/internal/handler"
	"github.com/versoit/diploma/be/logistics/internal/repository"
	"github.com/versoit/diploma/be/logistics/internal/usecase"
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
		repository.NewDeliveryRepository,
		repository.NewCourierRepository,
		usecase.NewLogisticsUseCase,
		handler.NewLogisticsHandler,
		handler.NewNatsHandler,
	),
	fx.Invoke(
		func(*sdktrace.TracerProvider) {},
		func(lc fx.Lifecycle, natsHandler *handler.NatsHandler, log *slog.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					log.Info("Starting NATS consumer")
					return natsHandler.Start()
				},
			})
		},
	),
)
