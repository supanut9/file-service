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

	// Initialize DB and get the instance
	dbConn, err := db.InitDB(&cfg.DB)
	if err != nil {
		log.Fatalf("❌ Could not initialize database: %v", err)
	}

	// Initialize R2 and get the client instance
	r2Client, err := config.NewR2Client()
	if err != nil {
		log.Fatalf("❌ Could not initialize R2 client: %v", err)
	}

	// Pass the instances to the router setup
	route.Setup(app, dbConn, r2Client)

	fmt.Println("App running on port:", cfg.Port)
	app.Listen(":" + cfg.Port)
}
