package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func NewRequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		rid := c.Get("X-Request-ID")
		if rid == "" {
			rid = uuid.New().String()
		}
		c.Set("X-Request-ID", rid)
		return c.Next()
	}
}
