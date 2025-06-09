package db

import (
	"fmt"
	"log"

	"github.com/supanut9/file-service/internal/config"
	"github.com/supanut9/file-service/internal/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(cfg *config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	log.Println("✅ Connected to MySQL database")

	err = migrate(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&entity.File{},
	)
	if err != nil {
		return fmt.Errorf("failed to run DB migrations: %v", err)
	}
	log.Println("✅ Database migrated successfully")
	return nil
}
