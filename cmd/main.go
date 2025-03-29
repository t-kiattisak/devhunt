package main

import (
	config "devhunt/configs"
	"devhunt/internal/delivery"
	"devhunt/internal/infrastructure"
	"devhunt/internal/repository"
	"devhunt/internal/usecase"
	"devhunt/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.Load()
	logger.Init()

	db := infrastructure.NewPostgresDB()

	toolRepo := repository.NewToolRepository(db)
	toolUsecase := usecase.NewToolUsecase(toolRepo)

	app := fiber.New()
	delivery.NewToolHandler(app, toolUsecase)

	app.Listen(":3000")
}
