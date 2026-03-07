package handlers

import (
	"github.com/gofiber/fiber/v2"
	"jodhere.backend/utils"
)

// HealthCheck godoc
// @Summary Health check
// @Description Return API health status
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func HealthCheck(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, "API is running", fiber.Map{
		"database": "connected",
		"version":  "1.0.0",
	})
}

// Welcome godoc
// @Summary Welcome
// @Description Returns a small index of API endpoints
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1 [get]
func Welcome(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, "Welcome to JodHere API", fiber.Map{
		"endpoints": []fiber.Map{
			{"method": "GET", "path": "/health", "description": "Health check"},
			{"method": "GET", "path": "/api/v1/auth", "description": "List all auth records"},
			{"method": "POST", "path": "/api/v1/auth", "description": "Create a new auth"},
			{"method": "GET", "path": "/api/v1/auth/:id", "description": "Get auth by ID"},
			{"method": "PUT", "path": "/api/v1/auth/:id", "description": "Update auth"},
			{"method": "DELETE", "path": "/api/v1/auth/:id", "description": "Delete auth"},

			{"method": "GET", "path": "/api/v1/places", "description": "List all places"},
			{"method": "POST", "path": "/api/v1/places", "description": "Create a new place"},
			{"method": "GET", "path": "/api/v1/places/:id", "description": "Get place by ID"},
			{"method": "PUT", "path": "/api/v1/places/:id", "description": "Update place"},
			{"method": "DELETE", "path": "/api/v1/places/:id", "description": "Delete place"},

			{"method": "GET", "path": "/api/v1/parking-zones", "description": "List all parking zones"},
			{"method": "POST", "path": "/api/v1/parking-zones", "description": "Create a new parking zone"},
			{"method": "GET", "path": "/api/v1/parking-zones/:id", "description": "Get parking zone by ID"},
			{"method": "PUT", "path": "/api/v1/parking-zones/:id", "description": "Update parking zone"},
			{"method": "DELETE", "path": "/api/v1/parking-zones/:id", "description": "Delete parking zone"},
		},
	})
}
