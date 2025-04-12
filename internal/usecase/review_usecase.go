package usecase

import (
	"context"
	"devhunt/internal/domain"
	"devhunt/internal/infrastructure"
	"devhunt/internal/repository"
	"fmt"
)

type ReviewUsecase struct {
	repo repository.ReviewRepository
}

func NewReviewUsecase(repo repository.ReviewRepository) *ReviewUsecase {
	return &ReviewUsecase{repo: repo}
}

func (u *ReviewUsecase) CreateReview(toolID int, userID string, rating int, comment string) error {
	if rating < 1 || rating > 5 {
		return fmt.Errorf("rating must be between 1 and 5")
	}

	err := u.repo.CreateReview(toolID, userID, rating, comment)
	if err != nil {
		return err
	}

	ctx := context.Background()
	cacheKeys := []string{
		fmt.Sprintf("tool:%d", toolID),
		fmt.Sprintf("rating:tool:%d", toolID),
		fmt.Sprintf("reviews:tool:%d", toolID),
	}
	for _, key := range cacheKeys {
		_ = infrastructure.Redis.Del(ctx, key).Err()
	}

	return nil
}

func (u *ReviewUsecase) GetReviews(toolID, limit, offset int) ([]domain.Review, error) {
	return u.repo.GetReviews(toolID, limit, offset)
}
