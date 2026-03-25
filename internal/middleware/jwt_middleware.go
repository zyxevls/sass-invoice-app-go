package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/zyxevls/internal/helpers"
)

func JWTMiddleware(c fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	token = strings.TrimPrefix(token, "Bearer ")

	_, err := helpers.ValidateToken(token)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return c.Next()
}
