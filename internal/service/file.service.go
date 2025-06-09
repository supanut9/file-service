package service

import (
	"context"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/supanut9/file-service/internal/config"
	"github.com/supanut9/file-service/internal/entity"
	"github.com/supanut9/file-service/internal/repository"
	"github.com/supanut9/file-service/internal/storage"
)

type FileService interface {
	UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, bucketName string, folderPath string, isPublic bool) (string, error)
}

type fileService struct {
	repo     repository.FileRepository
	r2Client *s3.Client
	r2Config config.R2Config
}

func NewFileService(repo repository.FileRepository, r2Client *s3.Client, r2Config config.R2Config) FileService {
	return &fileService{
		repo:     repo,
		r2Client: r2Client,
		r2Config: r2Config,
	}
}

func (s *fileService) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, bucketName string, folderPath string, isPublic bool) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// If no bucket name is provided in the request, use the default from config
	if bucketName == "" {
		bucketName = s.r2Config.BucketName
	}

	// Pass the r2Config to the storage function
	url, key, err := storage.UploadToR2(s.r2Client, s.r2Config, file, fileHeader, bucketName, folderPath, isPublic)
	if err != nil {
		return "", err
	}

	err = s.repo.Create(&entity.File{
		URL:   url,
		Title: fileHeader.Filename,
	})
	if err != nil {
		go storage.DeleteFromR2(s.r2Client, s.r2Config, bucketName, key)
		return "", err
	}

	return url, nil
}
