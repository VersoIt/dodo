package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type HealthHandler struct {
	log *zap.Logger
}

func NewHealthHandler(log *zap.Logger) *HealthHandler {
	return &HealthHandler{log: log}
}

func (h *HealthHandler) Check(c *fiber.Ctx) error {
	h.log.Debug("Health check requested")
	return c.JSON(fiber.Map{
		"status": "ok",
		"service": "gateway",
	})
}
