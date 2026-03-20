package http

import "github.com/gofiber/fiber/v3"

func NewRouter(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "API v1 is running",
		})
	})
}
