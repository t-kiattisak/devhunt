package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type VoteRepository interface {
	CreateVote(tollID int, userID string) error
}

type voteRepository struct {
	DB *pgxpool.Pool
}

func NewVoteRepository(db *pgxpool.Pool) VoteRepository {
	return &voteRepository{DB: db}
}

// CreateVote implements VoteRepository.
func (r *voteRepository) CreateVote(tollID int, userID string) error {

	ctx := context.Background()
	_, err := r.DB.Exec(ctx, `
		INSERT INTO tool_votes (tool_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (tool_id, user_id) DO NOTHING
	`, tollID, userID)
	return err
}
