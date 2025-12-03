package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/myplate/backend/internal/services"
)

func AuthMiddleware(authService *services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Требуется заголовок авторизации"})
		}
		
		// Remove "Bearer " prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
		
		userID, role, err := authService.ValidateJWT(token)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Неверный токен"})
		}
		
		c.Locals("user_id", userID)
		c.Locals("user_role", role)
		return c.Next()
	}
}

// AdminMiddleware проверяет, что пользователь является администратором
func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("user_role").(string)
		if !ok || role != "admin" {
			return c.Status(403).JSON(fiber.Map{"error": "Доступ запрещен. Требуются права администратора"})
		}
		return c.Next()
	}
}


