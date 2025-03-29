package delivery

import (
	"devhunt/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type ToolHandler struct {
	usecase *usecase.ToolUsecase
}

func NewToolHandler(app *fiber.App, usecase *usecase.ToolUsecase) {
	handler := &ToolHandler{usecase: usecase}
	app.Get("/tools", handler.GetAllTools)
}

func (h *ToolHandler) GetAllTools(c *fiber.Ctx) error {
	tools, err := h.usecase.GetAllTools()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching tools",
		})
	}
	return c.JSON(tools)
}
