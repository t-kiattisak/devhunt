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
