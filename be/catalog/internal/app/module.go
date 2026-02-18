package app

import (
	"context"

	"github.com/versoit/diploma/be/catalog/internal/config"
	"github.com/versoit/diploma/be/catalog/internal/handler"
	"github.com/versoit/diploma/be/catalog/internal/repository"
	"github.com/versoit/diploma/be/catalog/internal/usecase"
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
		repository.NewProductRepository,
		usecase.NewCatalogUseCase,
		handler.NewCatalogHandler,
	),
	fx.Invoke(func(*sdktrace.TracerProvider) {}),
)
