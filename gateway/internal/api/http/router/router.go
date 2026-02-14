package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/versoit/diploma/gateway/internal/api/http/handlers"
	"github.com/versoit/diploma/gateway/internal/api/http/middleware"
	"log/slog"
)

func SetupRoutes(
	app *fiber.App,
	log *slog.Logger,
	healthHandler *handlers.HealthHandler,
	authHandler *handlers.AuthHandler,
	catalogHandler *handlers.CatalogHandler,
	orderHandler *handlers.OrderHandler,
	kitchenHandler *handlers.KitchenHandler,
	logisticsHandler *handlers.LogisticsHandler,
	analyticsHandler *handlers.AnalyticsHandler,
) {
	app.Use(middleware.NewRequestIDMiddleware())
	app.Use(middleware.NewLoggerMiddleware(log))
	
	api := app.Group("/api/v1")
	// Health
	api.Get("/health", healthHandler.Check)

	// Auth
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Catalog (Public)
	catalog := api.Group("/catalog")
	catalog.Get("/products", catalogHandler.ListProducts)
	catalog.Get("/products/:id", catalogHandler.GetProduct)

	// Protected routes
	protected := api.Group("/", middleware.NewAuthMiddleware())
	protected.Get("/auth/me", authHandler.GetMe)

	// Orders (Client flow)
	orders := protected.Group("/orders")
	orders.Post("/", orderHandler.CreateOrder)
	orders.Get("/my", orderHandler.ListOrders)
	orders.Get("/:id", orderHandler.GetOrderStatus)
	orders.Post("/:id/pay", orderHandler.PayOrder)

	// Kitchen (Internal flow for cooks)
	kitchen := protected.Group("/kitchen")
	kitchen.Patch("/tickets/:id/status", kitchenHandler.UpdateStatus)
	
	// Logistics (Courier flow)
	logistics := protected.Group("/logistics")
	logistics.Post("/orders/:id/assign", logisticsHandler.AssignCourier)
	logistics.Post("/orders/:id/location", logisticsHandler.UpdateLocation)

	// Analytics
	analytics := protected.Group("/analytics")
	analytics.Get("/manager/:id", analyticsHandler.GetManagerKPI)
}
