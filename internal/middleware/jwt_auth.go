package middleware

import (
	"devhunt/pkg/jwt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		userID, err := jwt.ParseToken(tokenStr)
		if err != nil || userID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		c.Locals("user_id", userID)
		return c.Next()
	}
}
