package routes

import (
	"github.com/gofiber/fiber/v2"
	"jodhere.backend/handlers"
	"jodhere.backend/middleware"
)

func SetupRoutes(app *fiber.App) {
	// Swagger UI is served by the app-level swagger handler at /swagger/*

	app.Get("/health", handlers.HealthCheck)
	api := app.Group("/api/v1", middleware.RequireAuth)
	api.Get("/", handlers.Welcome)

	// register grouped routes in separate files
	RegisterProfileRoutes(api.Group("/profile"))
	RegisterPlaceRoutes(api.Group("/places"))
	RegisterZoneRoutes(api.Group("/parking-zones"))

	// catch-all handled by NotFound above
}
