package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/supanut9/file-service/db"
	"github.com/supanut9/file-service/internal/config"
	"github.com/supanut9/file-service/internal/route"
)

func main() {
	cfg := config.Load()
	app := fiber.New()

	db.InitDB(&cfg.DB)

	config.InitR2()

	route.Setup(app)

	fmt.Println("App running on port:", cfg.Port)
	app.Listen(":" + cfg.Port)
}
