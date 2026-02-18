package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/versoit/diploma/be/treasury/internal/config"
	"github.com/versoit/diploma/be/treasury/internal/app"
	"github.com/versoit/diploma/be/treasury/internal/handler"
	"github.com/versoit/diploma/pkg/common"
	"go.uber.org/fx"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	stdgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	fx.New(
		app.Module,
		fx.Provide(
			func(cfg *config.Config) *slog.Logger {
				return common.NewLogger(cfg.App.Name)
			},
		),
		fx.Invoke(RunServer),
	).Run()
}

func RunServer(lc fx.Lifecycle, handler *handler.TreasuryHandler, logger *slog.Logger, cfg *config.Config) {
	server := stdgrpc.NewServer(
		stdgrpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	handler.Register(server)
	reflection.Register(server)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := fmt.Sprintf(":%s", cfg.GRPC.Port)
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				return err
			}
			logger.Info("Starting gRPC server", slog.String("port", cfg.GRPC.Port))
			go func() {
				if err := server.Serve(lis); err != nil {
					logger.Error("gRPC server failed", slog.Any("error", err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping gRPC server")
			server.GracefulStop()
			return nil
		},
	})
}
