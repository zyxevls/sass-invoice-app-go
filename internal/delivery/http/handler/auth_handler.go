package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/zyxevls/internal/usecase"
)

type AuthHandler struct {
	usecase usecase.AuthUsecase
}

func NewAuthHandler(u usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{u}
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req struct {
		Email    string
		Password string
	}

	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	token, err := h.usecase.Login(req.Email, req.Password)
	if err != nil {
		return c.SendStatus(401)
	}

	return c.JSON(fiber.Map{"token": token})
}
