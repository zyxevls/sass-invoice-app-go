package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/zyxevls/internal/usecase"
)

type PaymentHandler struct {
	usecase usecase.PaymentUseCase
}

func NewPaymentHandler(u usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{u}
}

func (h *PaymentHandler) Create(c fiber.Ctx) error {
	type req struct {
		InvoiceID string `json:"invoice_id"`
		Email     string `json:"email"`
		Amount    int64  `json:"amount"`
	}

	var r req
	if err := c.Bind().Body(&r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	url, err := h.usecase.CreatePayment(r.InvoiceID, r.Email, r.Amount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"payment_url": url})
}

func (h *PaymentHandler) WebHook(c fiber.Ctx) error {
	type payload struct {
		OrderID           string `json:"order_id"`
		TransactionStatus string `json:"transaction_status"`
	}

	var p payload
	if err := c.Bind().Body(&p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.usecase.HandleWebHook(p.OrderID, p.TransactionStatus)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "webhook processed", "order_id": p.OrderID, "transaction_status": p.TransactionStatus})
}
