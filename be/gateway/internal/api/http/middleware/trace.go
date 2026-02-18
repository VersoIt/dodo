package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
)

func NewTraceMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		span := trace.SpanFromContext(c.UserContext())
		if span.SpanContext().IsValid() {
			traceID := span.SpanContext().TraceID().String()
			c.Set("X-Trace-Id", traceID)
		}
		return c.Next()
	}
}
