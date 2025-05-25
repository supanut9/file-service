package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/supanut9/file-service/internal/handler"
	"github.com/supanut9/file-service/internal/repository"
	"github.com/supanut9/file-service/internal/service"
)

func Setup(app *fiber.App) {
	// Repository
	fileRepo := repository.NewFileRepository()

	// Service
	fileService := service.NewFileService(fileRepo)

	// Handler
	fileHandler := handler.NewFileHandler(fileService)

	api := app.Group("/api/v1")

	fileHandler.RegisterFileRoutes(api.Group("/files"))

}
