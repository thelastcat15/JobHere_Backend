package routes

import (
	"github.com/gofiber/fiber/v2"
	"jobhere.backend/handlers"
)

// RegisterZoneRoutes registers parking-zone-related endpoints
func RegisterZoneRoutes(r fiber.Router) {
	r.Post("/", handlers.CreateParkingZone)
	r.Get("/", handlers.ListParkingZones)
	r.Get(":id", handlers.GetParkingZone)
	r.Put(":id", handlers.UpdateParkingZone)
	r.Delete(":id", handlers.DeleteParkingZone)
}
