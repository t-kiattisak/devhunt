package repository

import (
	"context"
	"devhunt/internal/domain"
	"devhunt/pkg/logger"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ToolRepository interface {
	GetAllTools() ([]domain.Tool, error)
}

type toolRepository struct {
	DB *pgxpool.Pool
}

func NewToolRepository(db *pgxpool.Pool) ToolRepository {
	return &toolRepository{DB: db}
}

func (r *toolRepository) GetAllTools() ([]domain.Tool, error) {
	rows, err := r.DB.Query(context.Background(), "SELECT id, name, description FROM tools")
	if err != nil {
		logger.Error("Error querying tools: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	var tools []domain.Tool
	for rows.Next() {
		var tool domain.Tool
		err := rows.Scan(&tool.ID, &tool.Name, &tool.Description)
		if err != nil {
			return nil, err
		}
		tools = append(tools, tool)
	}

	return tools, nil
}
