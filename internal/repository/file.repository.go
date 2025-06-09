package repository

import (
	"github.com/supanut9/file-service/internal/entity"
	"gorm.io/gorm"
)

type FileRepository interface {
	Create(file *entity.File) error
}

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) Create(file *entity.File) error {
	return r.db.Create(file).Error
}
