package repository

import (
	"context"
	"devhunt/internal/domain"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ToolRepository interface {
	GetAllTools() ([]domain.Tool, error)
	GetTools(search string, limit, offset int) ([]domain.Tool, error)
	GetToolsCursor(cursorID int, limit int) ([]domain.Tool, error)
	GetToolsCursorWithSearch(search string, cursorID int, limit int) ([]domain.Tool, error)
	GetToolByID(toolID int) (*domain.Tool, error)
	CountVotes(toolID int) (int, error)
	CountReviews(toolID int) (int, error)
	AvgRating(toolID int) (float64, error)
	GetTopTrending(by string, limit int) ([]domain.Tool, error)
}

type toolRepository struct {
	DB *pgxpool.Pool
}

func NewToolRepository(db *pgxpool.Pool) ToolRepository {
	return &toolRepository{DB: db}
}

func (r *toolRepository) GetAllTools() ([]domain.Tool, error) {
	return r.GetTools("", 10, 0)
}

func (r *toolRepository) GetTools(search string, limit, offset int) ([]domain.Tool, error) {
	ctx := context.Background()

	var rows pgx.Rows
	var err error
	var tools []domain.Tool

	if search == "" {
		query := `
			SELECT id, name, description
			FROM tools
			ORDER BY id DESC
			LIMIT $1 OFFSET $2
		`
		rows, err = r.DB.Query(ctx, query, limit, offset)
	} else {
		query := `
			SELECT id, name, description
			FROM tools
			WHERE name ILIKE $1
			ORDER BY id DESC
			LIMIT $2 OFFSET $3
		`
		searchParam := "%" + search + "%"
		rows, err = r.DB.Query(ctx, query, searchParam, limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tool domain.Tool
		if err := rows.Scan(&tool.ID, &tool.Name, &tool.Description); err != nil {
			return nil, err
		}
		tools = append(tools, tool)
	}

	return tools, nil
}

func (r *toolRepository) GetToolsCursor(cursorID int, limit int) ([]domain.Tool, error) {
	ctx := context.Background()
	query := `
		SELECT id, name, description
		FROM tools
		WHERE ($1 = 0 OR id < $1)
		ORDER BY id DESC
		LIMIT $2
	`
	rows, err := r.DB.Query(ctx, query, cursorID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tools []domain.Tool

	for rows.Next() {
		var tool domain.Tool
		if err := rows.Scan(&tool.ID, &tool.Name, &tool.Description); err != nil {
			return nil, err
		}
		tools = append(tools, tool)
	}
	return tools, nil
}

func (r *toolRepository) GetToolsCursorWithSearch(search string, cursorID int, limit int) ([]domain.Tool, error) {
	ctx := context.Background()
	query := `
		SELECT id, name, description
		FROM tools
		WHERE ($1 = '' OR name ILINK $1)
		AND ($2 = 0 OR id < $2)
		ORDER BY id DESC
		LIMIT $3
	`
	searchParam := "%"
	if search != "" {
		searchParam = "%" + search + "%"
	}

	rows, err := r.DB.Query(ctx, query, searchParam, cursorID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tools []domain.Tool
	for rows.Next() {
		var tool domain.Tool
		if err := rows.Scan(&tool.ID, &tool.Name, &tool.Description); err != nil {
			return nil, err
		}
		tools = append(tools, tool)
	}
	return tools, nil
}

// GetToolByID implements ToolRepository.
func (r *toolRepository) GetToolByID(toolID int) (*domain.Tool, error) {
	ctx := context.Background()

	query := `
		SELECT id, name, description
		FROM tools
		WHERE (id = $1)
	`
	row := r.DB.QueryRow(ctx, query, toolID)
	var tool domain.Tool
	err := row.Scan(&tool.ID, &tool.Name, &tool.Description)
	if err != nil {
		return nil, err
	}
	return &tool, nil
}

func (r *toolRepository) CountVotes(toolID int) (int, error) {
	ctx := context.Background()

	query := `SELECT COUNT(*) FROM tool_votes WHERE tool_id = $1`
	row := r.DB.QueryRow(ctx, query, toolID)

	var count int
	err := row.Scan(&count)
	return count, err
}

func (r *toolRepository) CountReviews(toolID int) (int, error) {
	ctx := context.Background()
	row := r.DB.QueryRow(ctx, `SELECT COUNT(*) FROM tool_review WHERE tool_id = $1`, toolID)
	var count int
	err := row.Scan(count)
	return count, err
}

func (r *toolRepository) AvgRating(toolID int) (float64, error) {
	ctx := context.Background()
	row := r.DB.QueryRow(ctx, `SELECT COALESCE(AVG(rating), 0)  FROM tool_reviews WHERE tool_id = $1`, toolID)
	var avg float64
	err := row.Scan(&avg)
	return avg, err
}

// GetTopTrending implements ToolRepository.
func (r *toolRepository) GetTopTrending(by string, limit int) ([]domain.Tool, error) {
	ctx := context.Background()

	var query string
	switch by {
	case "rating":
		query = `
			SELECT t.id, t.name, t.description
			FROM tools t
			LEFT JOIN tool_reviews r ON t.id = r.tool_id
			GROUP BY t.id
			ORDER BY COALESCE(AVG(r.rating), 0) DESC
			LIMIT $1
		`
	case "reviews":
		query = `
			SELECT t.id, t.name, t.description
			FROM tools t
			LEFT JOIN tool_reviews r ON t.id = r.tool_id
			GROUP BY t.id
			ORDER By COUNT(r.id) DESC
			LIMIT $1
		`
	default:
		query = `
			SELECT t.id, t.name, t.description
			FROM tools t
			LEFT JOIN tool_votes v ON t.id = v.tool_id
			GROUP BY t.id
			ORDER BY COUNT(v.id) DESC
			LIMIT $1
		`
	}

	rows, err := r.DB.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tools []domain.Tool
	for rows.Next() {
		var tool domain.Tool
		if err := rows.Scan(&tool.ID, &tool.Name, &tool.Description); err != nil {
			return nil, err
		}
		tools = append(tools, tool)
	}
	return tools, nil
}
