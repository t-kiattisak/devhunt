package usecase

import (
	"devhunt/internal/domain"
	"devhunt/internal/repository"
)

type ToolUsecase struct {
	repo repository.ToolRepository
}

func NewToolUsecase(repo repository.ToolRepository) *ToolUsecase {
	return &ToolUsecase{repo: repo}
}

func (u *ToolUsecase) GetAllTools() ([]domain.Tool, error) {
	return u.repo.GetAllTools()
}

func (u *ToolUsecase) GetTools(search string, limit int, offset int) ([]domain.Tool, error) {
	return u.repo.GetTools(search, limit, offset)
}

func (u *ToolUsecase) GetToolsCursor(cursorID int, limit int) ([]domain.Tool, error) {
	return u.repo.GetToolsCursor(cursorID, limit)
}

func (u *ToolUsecase) GetToolsCursorWithSearch(search string, cursorID int, limit int) ([]domain.Tool, error) {
	return u.repo.GetToolsCursorWithSearch(search, cursorID, limit)
}

type ToolDetail struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	VoteCount   int     `json:"vote_count"`
	ReviewCount int     `json:"review_count"`
	AvgRating   float64 `json:"avg_rating"`
}

func (u *ToolUsecase) GetToolByID(toolID int) (*ToolDetail, error) {
	tool, err := u.repo.GetToolByID(toolID)
	if err != nil {
		return nil, err
	}

	voteCount, err := u.repo.CountVotes(toolID)
	if err != nil {
		return nil, err
	}

	reviewCount, _ := u.repo.CountReviews(toolID)
	avgRating, _ := u.repo.AvgRating(toolID)

	return &ToolDetail{
		ID:          tool.ID,
		Name:        tool.Name,
		Description: tool.Description,
		VoteCount:   voteCount,
		ReviewCount: reviewCount,
		AvgRating:   avgRating,
	}, nil
}

func (u *ToolUsecase) GetTopTrending(by string, limit int) ([]domain.Tool, error) {
	return u.repo.GetTopTrending(by, limit)
}
