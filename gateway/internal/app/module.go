package app

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/versoit/diploma/gateway/internal/config"
	"github.com/versoit/diploma/gateway/internal/api/http/router"
	"github.com/versoit/diploma/gateway/internal/api/http/handlers"
	"github.com/versoit/diploma/gateway/internal/infrastructure/grpc"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(
	grpc.Module,
	fx.Provide(
		config.NewConfig,
		newLogger,
		newFiberApp,
		handlers.NewHealthHandler,
		handlers.NewAuthHandler,
		handlers.NewCatalogHandler,
		handlers.NewOrderHandler,
		handlers.NewKitchenHandler,
		handlers.NewLogisticsHandler,
	),
	fx.Invoke(
		router.SetupRoutes,
		startServer,
	),
)

func newLogger(cfg *config.Config) (*zap.Logger, error) {
	if cfg.LogLevel == "debug" {
		return zap.NewDevelopment()
	}
	return zap.NewProduction()
}

func newFiberApp() *fiber.App {
	return fiber.New(fiber.Config{
		AppName: "Diploma Gateway v1.0",
	})
}

func startServer(lc fx.Lifecycle, app *fiber.App, cfg *config.Config, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting gateway server", zap.String("port", cfg.AppPort))
			go func() {
				if err := app.Listen(":" + cfg.AppPort); err != nil {
					log.Fatal("Failed to start server", zap.Error(err))
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
