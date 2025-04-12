package repository

import (
	"context"
	"devhunt/internal/domain"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ReviewRepository interface {
	CreateReview(toolID int, userID string, rating int, comment string) error
	GetReviews(toolID, limit, offset int) ([]domain.Review, error)
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

// GetReviews implements ReviewRepository.
func (r *reviewRepository) GetReviews(toolID int, limit int, offset int) ([]domain.Review, error) {
	ctx := context.Background()
	query := `
		SELECT user_id, rating, comment, created_at
		FROM tool_reviews
		WHERE tool_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.DB.Query(ctx, query, toolID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []domain.Review
	for rows.Next() {
		var review domain.Review
		if err := rows.Scan(&review.UserID, &review.Rating, &review.Comment, &review.CreatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}
