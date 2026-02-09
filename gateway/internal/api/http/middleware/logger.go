package middleware

import (
	"time"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func NewLoggerMiddleware(log *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		
		err := c.Next()
		
		duration := time.Since(start)
		rid := c.Get("X-Request-ID")
		
		log.Info("HTTP Request",
			zap.String("request_id", rid),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("duration", duration),
			zap.String("ip", c.IP()),
		)
		
		return err
	}
}
