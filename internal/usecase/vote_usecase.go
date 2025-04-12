package usecase

import (
	"context"
	"devhunt/internal/infrastructure"
	"devhunt/internal/repository"
	"fmt"
)

type VoteUsecase struct {
	repo repository.VoteRepository
}

func NewVoteUsecase(repo repository.VoteRepository) *VoteUsecase {
	return &VoteUsecase{repo: repo}
}

func (u *VoteUsecase) Vote(toolID int, userID string) error {
	err := u.repo.CreateVote(toolID, userID)

	if err != nil {
		return err
	}

	ctx := context.Background()
	cacheKeys := []string{
		fmt.Sprintf("tool:%d", toolID),
		fmt.Sprintf("votes:tool:%d", toolID),
		fmt.Sprintf("votes:tool:%d", toolID),
		"top-trending:limit=10:by=votes",
		"top-trending:limit=5:by=votes",
	}

	for _, key := range cacheKeys {
		_ = infrastructure.Redis.Del(ctx, key).Err()
	}

	return nil
}
