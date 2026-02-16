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
	catalog.Post("/products", middleware.NewAuthMiddleware(), middleware.RequireRole("manager"), catalogHandler.CreateProduct)
	catalog.Put("/products/:id", middleware.NewAuthMiddleware(), middleware.RequireRole("manager"), catalogHandler.UpdateProduct)

	// Protected routes
	protected := api.Group("/", middleware.NewAuthMiddleware())
	protected.Get("/auth/me", authHandler.GetMe)
	protected.Patch("/auth/me", authHandler.UpdateProfile)

	// Orders (Client flow)
	orders := protected.Group("/orders")
	orders.Post("/", orderHandler.CreateOrder)
	orders.Get("/my", orderHandler.ListOrders)
	orders.Get("/all", middleware.RequireRole("chef", "courier", "manager"), orderHandler.ListAllOrders)
	orders.Get("/analytics", middleware.RequireRole("manager"), orderHandler.GetAnalytics)
	orders.Get("/:id", orderHandler.GetOrderStatus)
	orders.Patch("/:id/status", middleware.RequireRole("chef", "courier", "manager"), orderHandler.UpdateOrderStatus)
	orders.Post("/:id/pay", orderHandler.PayOrder)

	// Public/Client Promo check (Accessible by any authenticated user)
	protected.Get("/promos/check/:code", orderHandler.CheckPromoCode)

	// Promo Codes (Manager only)
	promos := protected.Group("/promos", middleware.RequireRole("manager"))
	promos.Get("/", orderHandler.ListPromos)
	promos.Post("/", orderHandler.CreatePromoCode)
	promos.Delete("/:id", orderHandler.DeletePromo)

	// Kitchen (Internal flow for cooks)
	kitchen := protected.Group("/kitchen", middleware.RequireRole("chef", "manager"))
	kitchen.Get("/tickets", kitchenHandler.ListTickets)
	kitchen.Patch("/tickets/:id/status", kitchenHandler.UpdateStatus)
	
	// Logistics (Courier flow)
	logistics := protected.Group("/logistics", middleware.RequireRole("courier", "manager"))
	logistics.Get("/deliveries", logisticsHandler.ListDeliveries)
	logistics.Patch("/orders/:id/status", logisticsHandler.UpdateStatus)
	logistics.Post("/orders/:id/assign", logisticsHandler.AssignCourier)
	logistics.Post("/orders/:id/location", logisticsHandler.UpdateLocation)
}
