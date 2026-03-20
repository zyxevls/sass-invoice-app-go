package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/zyxevls/internal/config"
	"github.com/zyxevls/internal/delivery/http"
	"github.com/zyxevls/internal/infrastructure/database"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	db := database.NewPostgres(cfg)
	defer db.Close()

	app := fiber.New()

	http.NewRouter(app)

	app.Listen(":8080")
}
