package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/supanut9/file-service/db"
	"github.com/supanut9/file-service/internal/config"
	"github.com/supanut9/file-service/internal/route"
)

func main() {
	cfg := config.Load()
	app := fiber.New()

	dbConn, err := db.InitDB(&cfg.DB)
	if err != nil {
		log.Fatalf("❌ Could not initialize database: %v", err)
	}

	r2Client, err := config.NewR2Client(cfg.R2)
	if err != nil {
		log.Fatalf("❌ Could not initialize R2 client: %v", err)
	}

	route.Setup(app, dbConn, r2Client, cfg.R2)

	fmt.Println("App running on port:", cfg.Port)
	app.Listen(":" + cfg.Port)
}
