package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"jodhere.backend/utils"
)

func RequireAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Authorization header missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization format")
	}

	tokenString := parts[1]

	claims, err := utils.ParseSupabaseJWT(tokenString)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired token")
	}

	// ใส่ข้อมูล user ลง context
	c.Locals("user_id", claims.Sub)
	c.Locals("email", claims.Email)

	return c.Next()
}
