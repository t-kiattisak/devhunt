package delivery

import (
	"context"
	"devhunt/internal/domain"
	"devhunt/internal/infrastructure"
	"devhunt/internal/usecase"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

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
	app.Get("/tools/top-trending", handler.GetTopTrending)
	app.Get("/tools/:id", handler.GetToolByID)
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
	cacheKey := fmt.Sprintf("tools:cursor=%d;limit=%d:search=%s", cursor, limit, search)

	ctx := context.Background()
	val, err := infrastructure.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var cached []domain.Tool
		if err := json.Unmarshal([]byte(val), &cached); err == nil {
			return c.JSON(cached)
		}
	}

	tools, err := h.usecase.GetToolsCursorWithSearch(search, cursor, limit)
	fmt.Print(err)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get tools",
		})
	}

	bytes, _ := json.Marshal(tools)
	_ = infrastructure.Redis.Set(ctx, cacheKey, bytes, time.Minute*5).Err()

	return c.JSON(tools)
}

func (h *ToolHandler) GetToolByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	toolID, _ := strconv.Atoi(idParam)

	tool, err := h.usecase.GetToolByID(toolID)
	fmt.Print(err)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "tool not found",
		})
	}

	return c.JSON(tool)
}

func (h *ToolHandler) GetTopTrending(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	by := c.Query("by", "votes")

	cacheKey := fmt.Sprintf("top-trending:limit=%d:by=%s", limit, by)

	ctx := context.Background()
	val, err := infrastructure.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var cached []domain.Tool
		if err := json.Unmarshal([]byte(val), &cached); err == nil {
			return c.JSON(cached)
		}
	}

	tools, err := h.usecase.GetTopTrending(by, limit)
	fmt.Print(err)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get trending tools",
		})
	}

	bytes, _ := json.Marshal(tools)
	_ = infrastructure.Redis.Set(ctx, cacheKey, bytes, time.Minute*5).Err()
	return c.JSON(tools)
}
