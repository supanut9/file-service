package service

import (
	"context"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/service/s3"
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
}

func NewFileService(repo repository.FileRepository, r2Client *s3.Client) FileService {
	return &fileService{
		repo:     repo,
		r2Client: r2Client,
	}
}

func (s *fileService) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, bucketName string, folderPath string, isPublic bool) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	url, key, err := storage.UploadToR2(s.r2Client, file, fileHeader, bucketName, folderPath, isPublic)
	if err != nil {
		return "", err
	}

	err = s.repo.Create(&entity.File{
		URL:   url,
		Title: fileHeader.Filename,
	})

	if err != nil {
		go storage.DeleteFromR2(s.r2Client, bucketName, key)

		return "", err
	}

	return url, nil
}
