package delivery

import (
	"devhunt/internal/usecase"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ReviewHandler struct {
	usecase *usecase.ReviewUsecase
}

func NewReviewHandler(app *fiber.App, u *usecase.ReviewUsecase) {
	handler := &ReviewHandler{usecase: u}
	app.Post("/tools/:id/reviews", handler.CreateReview)
	app.Get("/tools/:id/reviews", handler.GetReviews)
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

func (h *ReviewHandler) GetReviews(c *fiber.Ctx) error {
	toolID, _ := strconv.Atoi(c.Params("id"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	reviews, err := h.usecase.GetReviews(toolID, limit, offset)
	fmt.Print(err)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get reviews",
		})
	}
	return c.JSON(reviews)
}
