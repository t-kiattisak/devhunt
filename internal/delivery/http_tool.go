package delivery

import (
	"devhunt/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ToolHandler struct {
	usecase *usecase.ToolUsecase
}

func NewToolHandler(app *fiber.App, usecase *usecase.ToolUsecase) {
	handler := &ToolHandler{usecase: usecase}
	app.Get("/tools", handler.GetTools)
	app.Get("/all-tools", handler.GetAllTools)
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

func (h *ToolHandler) GetTools(c *fiber.Ctx) error {
	search := c.Query("search", "")
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	tools, err := h.usecase.GetTools(search, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get tools",
		})
	}

	return c.JSON(tools)

}
