package app

import (
	"context"

	"github.com/versoit/diploma/be/auth/internal/config"
	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/pkg/telemetry"
	"github.com/versoit/diploma/be/auth/internal/handler"
	"github.com/versoit/diploma/be/auth/internal/repository"
	"github.com/versoit/diploma/be/auth/internal/usecase"
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
		repository.NewUserRepository,
		usecase.NewAuthUseCase,
		handler.NewAuthHandler,
	),
	fx.Invoke(func(*sdktrace.TracerProvider) {}),
)
