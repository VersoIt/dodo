package main

import (
	"context"
	"log/slog"
	"net"

	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/be/kitchen/internal/app"
	"github.com/versoit/diploma/be/kitchen/internal/handler"
	"go.uber.org/fx"
	stdgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	fx.New(
		fx.Provide(
			func() *slog.Logger {
				return common.NewLogger("kitchen-service")
			},
		),
		app.Module,
		fx.Invoke(RunServer),
	).Run()
}

func RunServer(lc fx.Lifecycle, handler *handler.KitchenHandler, logger *slog.Logger) {
	server := stdgrpc.NewServer()
	handler.Register(server)
	reflection.Register(server)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", ":8080")
			if err != nil {
				return err
			}
			logger.Info("Starting gRPC server", slog.String("port", "8080"))
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
