package main

import (
	config "devhunt/configs"
	"devhunt/internal/delivery"
	"devhunt/internal/infrastructure"
	"devhunt/internal/repository"
	"devhunt/internal/usecase"
	"devhunt/pkg/logger"
	"devhunt/pkg/seeder"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.Load()
	logger.Init()

	db := infrastructure.NewPostgresDB()

	if os.Getenv("SEED") == "true" {
		seeder.SeedTools(db, 10000)
	}

	toolRepo := repository.NewToolRepository(db)
	toolUsecase := usecase.NewToolUsecase(toolRepo)

	app := fiber.New()
	delivery.NewToolHandler(app, toolUsecase)

	app.Listen(":3000")
}
