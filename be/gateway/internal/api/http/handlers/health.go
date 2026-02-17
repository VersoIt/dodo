package handlers

import (
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

type HealthHandler struct {
	log *slog.Logger
}

func NewHealthHandler(log *slog.Logger) *HealthHandler {
	return &HealthHandler{log: log}
}

func (h *HealthHandler) Check(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
	})
}
