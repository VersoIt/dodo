package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/nats-io/nats.go"
	"github.com/versoit/diploma/be/chat/internal/config"
	"github.com/versoit/diploma/be/chat/internal/handler"
	"github.com/versoit/diploma/be/chat/internal/repository"
	"github.com/versoit/diploma/be/chat/internal/service"
	"github.com/versoit/diploma/pkg/common"
	"github.com/versoit/diploma/pkg/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var Module = fx.Options(
	fx.Provide(
		func() (context.Context, context.CancelFunc) {
			return context.WithCancel(context.Background())
		},
		config.NewConfig,
		func(cfg *config.Config) string {
			return cfg.Database.URL
		},
		func(cfg *config.Config) *slog.Logger {
			return common.NewLogger(cfg.App.Name)
		},
		common.NewPGXPool,
		common.NewTransactionManager,
		telemetry.NewTracerProvider,

		NewNatsConn,
		service.NewHub,
		repository.NewMessageRepository,
		service.NewMessageService,
		handler.NewChatHandler,
	),
	fx.Invoke(
		startHub,
		startGRPCServer,
	),
)

func startHub(hub *service.Hub) {
	go hub.Run()
}

func startGRPCServer(lc fx.Lifecycle, h *handler.ChatHandler, cfg *config.Config, log *slog.Logger) {
	server := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	h.Register(server)
	reflection.Register(server)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := fmt.Sprintf(":%s", cfg.GRPC.Port)
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				return err
			}
			log.Info("Starting gRPC server", slog.String("port", cfg.GRPC.Port))
			go func() {
				if err := server.Serve(lis); err != nil {
					log.Error("gRPC server failed", slog.Any("error", err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping gRPC server")
			server.GracefulStop()
			return nil
		},
	})
}

func NewNatsConn(cfg *config.Config, log *slog.Logger) (*nats.Conn, error) {
	nc, err := nats.Connect(cfg.Nats.URL)
	if err != nil {
		return nil, err
	}
	log.Info("Connected to NATS", slog.String("url", cfg.Nats.URL))
	return nc, nil
}
