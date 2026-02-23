package middleware

import "github.com/gofiber/fiber/v2"

// Simple placeholder auth middleware
func RequireAuth(c *fiber.Ctx) error {
	// TODO: implement authentication check
	return c.Next()
}
