package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/zyxevls/internal/delivery/http/handler"
	"github.com/zyxevls/internal/repository"
	"github.com/zyxevls/internal/usecase"
)

func NewRouter(app *fiber.App, db *sqlx.DB) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	//init layer
	invoiceRepo := repository.NewInvoiceRepository(db)
	invoiceUsecase := usecase.NewInvoiceUsecase(invoiceRepo)
	invoiceHandler := handler.NewInvoiceHandler(invoiceUsecase)

	v1.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "API v1 is running",
		})
	})

	v1.Post("/invoices", invoiceHandler.Create)
	v1.Get("/invoices", invoiceHandler.GetAll)
}
