package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/zyxevls/internal/usecase"
)

type DashboardHandler struct {
	usecase usecase.DashboardUsecase
}

func NewDashboardHandler(u usecase.DashboardUsecase) *DashboardHandler {
	return &DashboardHandler{u}
}

func (h *DashboardHandler) Get(c fiber.Ctx) error {
	data, err := h.usecase.GetDashboard()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   data,
	})
}

func (h *DashboardHandler) GetTopCustomer(c fiber.Ctx) error {
	data, err := h.usecase.GetTopCustomer()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   data,
	})
}

func (h *DashboardHandler) GetRecentTransaction(c fiber.Ctx) error {
	data, err := h.usecase.GetRecentTransaction()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   data,
	})
}
