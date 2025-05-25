package db

import (
	"fmt"
	"log"

	"github.com/supanut9/file-service/internal/config"
	"github.com/supanut9/file-service/internal/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.DBConfig) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	log.Println("✅ Connected to MySQL database")

	Migrate() // 👈 ACTUALLY CALL the migrate function
}

func Migrate() {
	err := DB.AutoMigrate(
		&entity.File{},
	)
	if err != nil {
		log.Fatalf("❌ Failed to run DB migrations: %v", err)
	}
	log.Println("✅ Database migrated successfully")
}
