package routes

import (
	"github.com/gofiber/fiber/v2"
	"jodhere.backend/handlers"
)

// RegisterPlaceRoutes registers place-related endpoints
func RegisterPlaceRoutes(r fiber.Router) {
	r.Get("/", handlers.ListPlaces)
	r.Get("/:place_id", handlers.GetPlace)
	r.Get("/:place_id/zones/:zone_id/slots", handlers.GetParkingSlots)
}
