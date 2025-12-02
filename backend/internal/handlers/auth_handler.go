package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/myplate/backend/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type TelegramAuthRequest struct {
	InitData string `json:"init_data"`
}

type AuthResponse struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}

func (h *AuthHandler) AuthenticateTelegram(c *fiber.Ctx) error {
	var req TelegramAuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверное тело запроса"})
	}
	
	user, token, err := h.authService.AuthenticateTelegram(req.InitData)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}
	
	return c.JSON(AuthResponse{
		User:  user,
		Token: token,
	})
}

// TestAuth - для разработки и тестирования, создает тестового пользователя и возвращает токен
func (h *AuthHandler) TestAuth(c *fiber.Ctx) error {
	user, token, err := h.authService.CreateTestUser()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
	return c.JSON(AuthResponse{
		User:  user,
		Token: token,
	})
}


