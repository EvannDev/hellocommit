package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/hellocommit/api/internal/repositories"
)

func Auth(userRepo *repositories.UserRepository) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		user, err := userRepo.GetByAccessToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		c.Locals("userID", user.ID)
		return c.Next()
	}
}
