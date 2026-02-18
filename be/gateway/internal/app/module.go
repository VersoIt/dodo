package app

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/gateway/internal/config"
	"github.com/versoit/diploma/gateway/internal/api/http/router"
	"github.com/versoit/diploma/gateway/internal/api/http/handlers"
	"github.com/versoit/diploma/gateway/internal/infrastructure/grpc"
	"github.com/versoit/diploma/pkg/telemetry"
	"go.uber.org/fx"
	"log/slog"
)

var Module = fx.Options(
	grpc.Module,
	fx.Provide(
		config.NewConfig,
		func(cfg *config.Config) *slog.Logger {
			return common.NewLogger(cfg.App.Name)
		},
		telemetry.NewTracerProvider,
		newFiberApp,
		handlers.NewHealthHandler,
		handlers.NewAuthHandler,
		handlers.NewCatalogHandler,
		handlers.NewOrderHandler,
		handlers.NewKitchenHandler,
		handlers.NewLogisticsHandler,
		handlers.NewChatHandler,
	),
	fx.Invoke(
		router.SetupRoutes,
		startServer,
	),
)

func newFiberApp() *fiber.App {
	return fiber.New(fiber.Config{
		AppName:      "Diploma Gateway v1.0",
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	})
}

func startServer(lc fx.Lifecycle, app *fiber.App, cfg *config.Config, log *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting gateway server", "port", cfg.HTTP.Port)
			go func() {
				if err := app.Listen(":" + cfg.HTTP.Port); err != nil {
					log.Error("Failed to start server", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping gateway server")
			return app.Shutdown()
		},
	})
}
