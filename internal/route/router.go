package route

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/supanut9/file-service/internal/handler"
	"github.com/supanut9/file-service/internal/repository"
	"github.com/supanut9/file-service/internal/service"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB, r2Client *s3.Client) {
	// Repository
	fileRepo := repository.NewFileRepository(db)

	// Service
	fileService := service.NewFileService(fileRepo, r2Client)

	// Handler
	fileHandler := handler.NewFileHandler(fileService)

	api := app.Group("/api/v1")

	fileHandler.RegisterFileRoutes(api.Group("/files"))
}
