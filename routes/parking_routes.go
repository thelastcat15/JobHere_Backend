package routes

import (
	"github.com/gofiber/fiber/v2"
	"jodhere.backend/handlers"
)

// RegisterParkingRoutes registers parking-related endpoints
func RegisterParkingRoutes(r fiber.Router) {
	r.Get("/", handlers.ListParking)
	r.Get("/:parking_id", handlers.GetParking)
	r.Get("/:parking_id/zones/:zone_id/slots", handlers.GetParkingSlots)
}
