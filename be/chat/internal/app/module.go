package app

import (
	"context"
	"log/slog"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"github.com/versoit/diploma/be/chat/internal/config"
	"github.com/versoit/diploma/be/chat/internal/handler"
	"github.com/versoit/diploma/be/chat/internal/repository"
	"github.com/versoit/diploma/be/chat/internal/service"
	"github.com/versoit/diploma/pkg/common"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		func() (context.Context, context.CancelFunc) {
			return context.WithCancel(context.Background())
		},
		config.NewConfig,
		func() *slog.Logger {
			return common.NewLogger("chat")
		},
		common.NewPGXPool,
		common.NewTransactionManager,

		NewNatsConn,
		service.NewHub,
		repository.NewMessageRepository,
		service.NewMessageService,
		handler.NewChatHandler,
		newFiberApp,
	),
	fx.Invoke(
		setupRoutes,
		startHub,
		startServer,
	),
)

func newFiberApp() *fiber.App {
	return fiber.New(fiber.Config{
		AppName: "Pizza Chat Service",
	})
}

func setupRoutes(app *fiber.App, h *handler.ChatHandler) {
	api := app.Group("/chat")
	api.Get("/history", h.GetHistory)
	api.Get("/sync", h.GetSync)

	app.Get("/ws/chat", h.WSUpgrade, websocket.New(h.WebSocketHandler))
}

func startHub(hub *service.Hub) {
	go hub.Run()
}

func startServer(lc fx.Lifecycle, app *fiber.App, cfg *config.Config, log *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting chat service server", "port", cfg.AppPort)
			go func() {
				if err := app.Listen(":" + cfg.AppPort); err != nil {
					log.Error("Failed to start server", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping chat service server")
			return app.Shutdown()
		},
	})
}

func NewNatsConn(cfg *config.Config, log *slog.Logger) (*nats.Conn, error) {
	nc, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		return nil, err
	}
	log.Info("Connected to NATS", slog.String("url", cfg.NatsURL))
	return nc, nil
}
