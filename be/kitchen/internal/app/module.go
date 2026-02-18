package app

import (
	"context"

	"github.com/versoit/diploma/be/kitchen/internal/config"
	"github.com/versoit/diploma/be/kitchen/internal/handler"
	"github.com/versoit/diploma/be/kitchen/internal/repository"
	"github.com/versoit/diploma/be/kitchen/internal/usecase"
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
		repository.NewTicketRepository,
		usecase.NewKitchenUseCase,
		handler.NewKitchenHandler,
	),
	fx.Invoke(func(*sdktrace.TracerProvider) {}),
)
