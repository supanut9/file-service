// handler/file_handler.go

package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/supanut9/file-service/internal/dto"
	"github.com/supanut9/file-service/internal/service"
	handler "github.com/supanut9/file-service/internal/utils"
)

var validate = validator.New()

type FileHandler struct {
	service service.FileService
}

func NewFileHandler(s service.FileService) *FileHandler {
	return &FileHandler{service: s}
}

func (h *FileHandler) RegisterFileRoutes(router fiber.Router) {
	router.Post("/", h.UploadFile)
}

func (h *FileHandler) UploadFile(c *fiber.Ctx) error {
	query := new(dto.UploadQueryDTO)

	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse query parameters.",
			"details": err.Error(),
		})
	}

	if err := validate.Struct(query); err != nil {
		validationErrors := handler.FormatValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed.",
			"details": validationErrors,
		})
	}

	bucketName := query.BucketName

	if query.IsPublic {
		bucketName = "public"
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Failed to read file from form.",
			"details": err.Error(),
		})
	}

	url, err := h.service.UploadFile(c.Context(), fileHeader, bucketName, query.FolderPath, query.IsPublic)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to process file.",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Upload successful.",
		"data": fiber.Map{
			"filename": fileHeader.Filename,
			"url":      url,
		},
	})
}
