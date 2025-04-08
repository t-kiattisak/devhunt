package delivery

import (
	"devhunt/internal/usecase"
	"fmt"
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
	app.Get("/tools/cursor", handler.GetToolsCursor)
	app.Get("/tools/cursor-search", handler.GetToolsCursor)
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

func (h *ToolHandler) GetToolsCursor(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	cursor, _ := strconv.Atoi(c.Query("cursor", "0"))
	tools, err := h.usecase.GetToolsCursor(cursor, limit)
	fmt.Print(err)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get tools",
		})
	}
	return c.JSON(tools)
}

func (h *ToolHandler) GetToolsCursorWithSearch(c *fiber.Ctx) error {
	search := c.Query("search", "")
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	cursor, _ := strconv.Atoi(c.Query("cursor", "0"))
	tools, err := h.usecase.GetToolsCursorWithSearch(search, cursor, limit)
	fmt.Print(err)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get tools",
		})
	}
	return c.JSON(tools)
}
