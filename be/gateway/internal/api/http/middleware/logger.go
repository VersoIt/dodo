package middleware

import (
	"time"
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

func NewLoggerMiddleware(log *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		
		err := c.Next()
		
		duration := time.Since(start)
		rid := c.Get("X-Request-ID")
		
		log.Info("HTTP Request",
			slog.String("request_id", rid),
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", c.Response().StatusCode()),
			slog.Duration("duration", duration),
			slog.String("ip", c.IP()),
		)
		
		return err
	}
}
