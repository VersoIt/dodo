package app

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/gateway/internal/config"
	"github.com/versoit/diploma/gateway/internal/api/http/router"
	"github.com/versoit/diploma/gateway/internal/api/http/handlers"
	"github.com/versoit/diploma/gateway/internal/infrastructure/grpc"
	"go.uber.org/fx"
	"log/slog"
)

var Module = fx.Options(
	grpc.Module,
	fx.Provide(
		config.NewConfig,
		func() *slog.Logger {
			return common.NewLogger("gateway")
		},
		newFiberApp,
		handlers.NewHealthHandler,
		handlers.NewAuthHandler,
		handlers.NewCatalogHandler,
		handlers.NewOrderHandler,
		handlers.NewKitchenHandler,
		handlers.NewLogisticsHandler,
		handlers.NewAnalyticsHandler,
	),
	fx.Invoke(
		router.SetupRoutes,
		startServer,
	),
)

func newFiberApp() *fiber.App {
	return fiber.New(fiber.Config{
		AppName: "Diploma Gateway v1.0",
	})
}

func startServer(lc fx.Lifecycle, app *fiber.App, cfg *config.Config, log *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting gateway server", "port", cfg.AppPort)
			go func() {
				if err := app.Listen(":" + cfg.AppPort); err != nil {
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
