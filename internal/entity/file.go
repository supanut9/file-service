package entity

import "time"

type File struct {
	ID        uint   `gorm:"primaryKey"`
	URL       string `gorm:"type:text;not null"`
	Title     string `gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
