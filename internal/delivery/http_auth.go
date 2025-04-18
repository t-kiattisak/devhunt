package delivery

import (
	"devhunt/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func NewAuthHandler(app *fiber.App) {
	app.Post("/v1/login", func(c *fiber.Ctx) error {
		var body struct {
			UserID string `json:"user_id"`
		}
		if err := c.BodyParser(&body); err != nil || body.UserID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}

		token, err := jwt.GenerateToken(body.UserID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "token failed"})
		}

		return c.JSON(fiber.Map{"token": token})
	})
}
