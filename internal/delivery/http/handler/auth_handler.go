package handler

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
	"github.com/zyxevls/internal/usecase"
)

type AuthHandler struct {
	usecase usecase.AuthUsecase
}

func NewAuthHandler(u usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{u}
}

func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Use manual JSON unmarshal for maximum compatibility (handles missing Content-Type)
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body", "details": err.Error()})
	}

	// Validate required fields
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name, email, and password are required"})
	}

	// Default role is always "user"
	role := "user"

	if err := h.usecase.Register(req.Name, req.Email, req.Password, role); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Failed to register user"})
	}

	return c.JSON(fiber.Map{"message": "User registered successfully", "email": req.Email})
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Use manual JSON unmarshal for maximum compatibility (handles missing Content-Type)
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body", "details": err.Error()})
	}

	token, err := h.usecase.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized", "details": "Invalid email or password"})
	}

	return c.JSON(fiber.Map{"message": "Login successful", "email": req.Email, "token": token})
}
