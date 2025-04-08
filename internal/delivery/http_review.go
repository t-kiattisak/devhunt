package delivery

import (
	"devhunt/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ReviewHandler struct {
	usecase *usecase.ReviewUsecase
}

func NewReviewHandler(app *fiber.App, u *usecase.ReviewUsecase) {
	handler := &ReviewHandler{usecase: u}
	app.Post("/tools/:id/reviews", handler.CreateReview)
}

type reviewRequest struct {
	UserID  string `json:"user_id"`
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

func (h *ReviewHandler) CreateReview(c *fiber.Ctx) error {
	idParam := c.Params("id")
	toolID, _ := strconv.Atoi(idParam)

	var req reviewRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	err := h.usecase.CreateReview(toolID, req.UserID, req.Rating, req.Comment)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusCreated)
}
