package routes

import (
	"github.com/gofiber/fiber/v2"
	"jobhere.backend/handlers"
)

func SetupRoutes(app *fiber.App) {
	// Swagger UI is served by the app-level swagger handler at /swagger/*

	app.Get("/health", handlers.HealthCheck)
	api := app.Group("/api/v1")
	api.Get("/", handlers.Welcome)

	// register grouped routes in separate files
	RegisterAuthRoutes(api.Group("/auth"))
	RegisterPlaceRoutes(api.Group("/places"))
	RegisterZoneRoutes(api.Group("/parking-zones"))

	// catch-all handled by NotFound above
}
