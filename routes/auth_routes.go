package routes

import (
	"github.com/gofiber/fiber/v2"
	"jobhere.backend/handlers"
)

// RegisterAuthRoutes registers auth-related endpoints
func RegisterAuthRoutes(r fiber.Router) {
	r.Post("/", handlers.CreateAuth)
	r.Get("/", handlers.ListAuth)
	r.Get(":id", handlers.GetAuth)
	r.Put(":id", handlers.UpdateAuth)
	r.Delete(":id", handlers.DeleteAuth)
}
