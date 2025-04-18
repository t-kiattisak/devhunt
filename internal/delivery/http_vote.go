package delivery

import (
	"devhunt/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type VoteHandler struct {
	usecase *usecase.VoteUsecase
}

func NewVoteHandler(app fiber.Router, u *usecase.VoteUsecase) {
	handler := &VoteHandler{usecase: u}
	app.Post("/tools/:id/vote", handler.Vote)
}

type voteRequest struct {
	UserID string `json:"user_id"`
}

func (h *VoteHandler) Vote(c *fiber.Ctx) error {
	idParam := c.Params("id")
	toolID, _ := strconv.Atoi(idParam)

	var req voteRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	err := h.usecase.Vote(toolID, req.UserID)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "already voted"})
	}

	return c.SendStatus(fiber.StatusCreated)
}
