package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ReviewRepository interface {
	CreateReview(toolID int, userID string, rating int, comment string) error
}

type reviewRepository struct {
	DB *pgxpool.Pool
}

func NewReviewRepository(db *pgxpool.Pool) ReviewRepository {
	return &reviewRepository{DB: db}
}

// CreateReview implements ReviewRepository.
func (r *reviewRepository) CreateReview(toolID int, userID string, rating int, comment string) error {
	ctx := context.Background()
	_, err := r.DB.Exec(ctx, `
		INSERT INTO tool_reviews (tool_id, user_id, rating, comment)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (tool_id, user_id) DO NOTHING
	`, toolID, userID, rating, comment)

	return err
}
