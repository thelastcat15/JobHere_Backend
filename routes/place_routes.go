package routes

import (
	"github.com/gofiber/fiber/v2"
	"jodhere.backend/handlers"
)

// RegisterPlaceRoutes registers place-related endpoints
func RegisterPlaceRoutes(r fiber.Router) {
	// r.Post("/", handlers.CreatePlace)
	r.Get("/", handlers.ListPlaces)
	r.Get(":id", handlers.GetPlace)
	// r.Put(":id", handlers.UpdatePlace)
	r.Delete(":id", handlers.DeletePlace)
}
