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
