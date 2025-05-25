package repository

import (
	"github.com/supanut9/file-service/db"
	"github.com/supanut9/file-service/internal/entity"
)

type FileRepository interface {
	Create(file *entity.File) error
}

type fileRepository struct{}

func NewFileRepository() FileRepository {
	return &fileRepository{}
}

func (r *fileRepository) Create(file *entity.File) error {
	return db.DB.Create(file).Error
}
