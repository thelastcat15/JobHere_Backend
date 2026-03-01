package routes

import (
	"github.com/gofiber/fiber/v2"
	"jodhere.backend/handlers"
)

// RegisterProfileRoutes registers profile-related endpoints
func RegisterProfileRoutes(r fiber.Router) {
	r.Get("/", handlers.GetProfile)
	r.Put("/", handlers.UpdateProfile)
}
