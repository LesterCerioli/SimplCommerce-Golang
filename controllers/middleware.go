package controllers

import (
	"github.com/gofiber/fiber/v3"
	"simplcommerce/services/implementations"
)

func AuthRequired(authService *implementations.AuthService) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" || len(token) < 7 || token[:7] != "Bearer " {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid authorization header"})
		}

		tokenString := token[7:]
		userID, email, roles, err := authService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		c.Locals("userID", userID)
		c.Locals("email", email)
		c.Locals("roles", roles)
		return c.Next()
	}
}

func AdminRequired() fiber.Handler {
	return func(c fiber.Ctx) error {
		roles, ok := c.Locals("roles").([]string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Admin access required"})
		}
		for _, role := range roles {
			if role == "admin" {
				return c.Next()
			}
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Admin access required"})
	}
}
