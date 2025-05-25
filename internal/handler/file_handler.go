package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/supanut9/file-service/internal/service"
)

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
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read file: " + err.Error(),
		})
	}

	bucketName := c.Query("bucketName")
	folderPath := c.Query("folderPath")
	isPublic := c.Query("isPublic", "true")
	public := isPublic == "true"

	url, err := h.service.UploadFile(c.Context(), fileHeader, bucketName, folderPath, public)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process file: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Upload successful",
		"filename": fileHeader.Filename,
		"url":      url,
	})
}
