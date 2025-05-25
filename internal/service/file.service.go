package service

import (
	"context"
	"mime/multipart"

	"github.com/supanut9/file-service/internal/entity"
	"github.com/supanut9/file-service/internal/repository"
	"github.com/supanut9/file-service/internal/storage"
)

type FileService interface {
	UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, bucketName string, folderPath string, isPublic bool) (string, error)
}

type fileService struct {
	repo repository.FileRepository
}

func NewFileService(repo repository.FileRepository) FileService {
	return &fileService{repo: repo}
}

func (s *fileService) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, bucketName string, folderPath string, isPublic bool) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	url, err := storage.UploadToR2(file, fileHeader.Filename, bucketName, folderPath, isPublic)
	if err != nil {
		return "", err
	}

	err = s.repo.Create(&entity.File{
		URL:   url,
		Title: fileHeader.Filename,
	})
	if err != nil {
		return "", err
	}

	return url, nil
}
