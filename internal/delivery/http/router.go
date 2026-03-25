package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/zyxevls/internal/config"
	"github.com/zyxevls/internal/delivery/http/handler"
	"github.com/zyxevls/internal/infrastructure/email"
	"github.com/zyxevls/internal/infrastructure/midtrans"
	"github.com/zyxevls/internal/infrastructure/pdf"
	"github.com/zyxevls/internal/infrastructure/redis"
	"github.com/zyxevls/internal/repository"
	"github.com/zyxevls/internal/usecase"
)

func NewRouter(app *fiber.App, db *sqlx.DB, cfg *config.Config) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	//init layer
	userRepo := repository.NewUserRepository(db)
	invoiceRepo := repository.NewInvoiceRepository(db)
	midtransSvc := midtrans.NewMidtransService(cfg)
	invoiceUsecase := usecase.NewInvoiceUsecase(invoiceRepo, cfg, midtransSvc)
	invoiceHandler := handler.NewInvoiceHandler(invoiceUsecase)
	paymentRepo := repository.NewPaymentRepository(db)
	paymentUsecase := usecase.NewPaymentUseCase(paymentRepo, invoiceRepo, midtransSvc, email.NewEmailService(cfg), pdf.NewPDFService())
	paymentHandler := handler.NewPaymentHandler(paymentUsecase)

	// Initialize Redis for caching
	rdb := redis.NewRedis()

	authUsecase := usecase.NewAuthUsecase(userRepo)
	authHandler := handler.NewAuthHandler(authUsecase)

	dashboardRepo := repository.NewDashboardRepositoryWithRedis(db, rdb)
	dashboardUsecase := usecase.NewDashboardUsecase(dashboardRepo)
	dashboardHandler := handler.NewDashboardHandler(dashboardUsecase)

	v1.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "API v1 is running",
		})
	})

	v1.Post("/auth/register", authHandler.Register)
	v1.Post("/auth/login", authHandler.Login)
	v1.Post("/invoices", invoiceHandler.Create)
	v1.Post("/payments", paymentHandler.Create)
	v1.Post("/payments/webhook", paymentHandler.WebHook)
	v1.Get("/invoices", invoiceHandler.GetAll)
	v1.Get("/dashboard", dashboardHandler.Get)
	v1.Get("/dashboard/top-customer", dashboardHandler.GetTopCustomer)
	v1.Get("/dashboard/recent-transaction", dashboardHandler.GetRecentTransaction)
}
