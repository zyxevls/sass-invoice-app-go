package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/zyxevls/internal/usecase"
)

type InvoiceHandler struct {
	usecase usecase.InvoiceUsecase
}

func NewInvoiceHandler(u usecase.InvoiceUsecase) *InvoiceHandler {
	return &InvoiceHandler{u}
}

func (h *InvoiceHandler) Create(c fiber.Ctx) error {
	var req usecase.CreateInvoiceRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	err := h.usecase.CreateInvoice(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "invoice created"})
}

func (h *InvoiceHandler) GetAll(c fiber.Ctx) error {
	data, err := h.usecase.GetInvoices()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(data)
}
