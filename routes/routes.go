package routes

import (
	"github.com/gofiber/fiber/v2"
	"jodhere.backend/handlers"
	"jodhere.backend/middleware"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/health", handlers.HealthCheck)
	api := app.Group("/api/v1", middleware.RequireAuth)
	api.Get("/", handlers.Welcome)

	RegisterProfileRoutes(api.Group("/profile"))
	RegisterPlaceRoutes(api.Group("/places"))

	// catch-all handled by NotFound above
}
