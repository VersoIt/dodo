package main

import (
	"context"
	"net"

	"github.com/versoit/diploma/services/auth/internal/api/grpc"
	"github.com/versoit/diploma/services/auth/internal/app"
	"go.uber.org/fx"
	"go.uber.org/zap"
	stdgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	fx.New(
		fx.Provide(
			zap.NewDevelopment,
		),
		app.Module,
		fx.Invoke(RunServer),
	).Run()
}

func RunServer(lc fx.Lifecycle, handler *grpc.AuthHandler, logger *zap.Logger) {
	server := stdgrpc.NewServer()
	handler.Register(server)
	reflection.Register(server)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", ":8080")
			if err != nil {
				return err
			}
			logger.Info("Starting gRPC server", zap.String("port", "8080"))
			go func() {
				if err := server.Serve(lis); err != nil {
					logger.Error("gRPC server failed", zap.Error(err))
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
