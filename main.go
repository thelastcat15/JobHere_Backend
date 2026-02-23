// @title Bakery POS API
// @version 1.0
// @description This is a Bakery POS API documentation
// @host localhost:5000
// @BasePath /api
// @schemes http

package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"jobhere.backend/config"
	_ "jobhere.backend/docs"
	"jobhere.backend/routes"
)

// @title JodHere API
// @version 1.0
// @description This is a JodHere API documentation
// @host localhost:5000
// @BasePath /api
// @schemes http
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

	// Create Fiber app with custom config
	app := fiber.New(fiber.Config{
		AppName:       "JobHere Backend v1.0.0",
		ErrorHandler:  errorHandler,
		StrictRouting: false,
	})

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path}\n",
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, http://localhost:5173",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Content-Type,Authorization",
		MaxAge:       3600,
	}))

	routes.SetupRoutes(app)

	// Serve swagger UI at /swagger/*
	app.Get("/swagger/*", swagger.HandlerDefault)

	port := ":3000"
	log.Printf("🚀 Starting JobHere Backend server on %s", port)

	if err := app.Listen(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// errorHandler is the global error handler for the Fiber app
func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"status":  code,
		"message": message,
		"error":   err.Error(),
	})
}
