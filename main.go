package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"jobhere.backend/config"
)

func main() {
	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	// Initialize Supabase database connection
	if err := config.InitSupabaseDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Close database connection on exit
	defer config.CloseDB()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "JobHere Backend v1.0.0",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, http://localhost:5173",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH",
		AllowHeaders: "Content-Type,Authorization",
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "JobHere API is running",
		})
	})

	// Routes placeholder
	api := app.Group("/api/v1")

	// Example routes (update as needed)
	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to JobHere API",
		})
	})

	// Start server
	log.Println("🚀 Starting JobHere Backend server on :3000")
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
