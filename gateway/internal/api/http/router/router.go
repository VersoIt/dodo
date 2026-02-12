package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/versoit/diploma/gateway/internal/api/http/handlers"
	"github.com/versoit/diploma/gateway/internal/api/http/middleware"
	"go.uber.org/zap"
)

func SetupRoutes(
	app *fiber.App,
	log *zap.Logger,
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

	// Catalog
	catalog := api.Group("/catalog")
	catalog.Get("/products", catalogHandler.ListProducts)
	catalog.Get("/products/:id", catalogHandler.GetProduct)

	// Orders (Client flow)
	orders := api.Group("/orders")
	orders.Post("/", orderHandler.CreateOrder)
	orders.Get("/:id", orderHandler.GetOrderStatus)
	orders.Post("/:id/pay", orderHandler.PayOrder)

	// Kitchen (Internal flow for cooks)
	kitchen := api.Group("/kitchen")
	kitchen.Patch("/tickets/:id/status", kitchenHandler.UpdateStatus)
	
	// Logistics (Courier flow)
	logistics := api.Group("/logistics")
	logistics.Post("/orders/:id/assign", logisticsHandler.AssignCourier)
	logistics.Post("/orders/:id/complete", logisticsHandler.CompleteDelivery)
	logistics.Post("/orders/:id/location", logisticsHandler.UpdateLocation)

	// Analytics
	analytics := api.Group("/analytics")
	analytics.Get("/manager/:id", analyticsHandler.GetManagerKPI)
}
