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

type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
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

// Register регистрирует нового пользователя
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверное тело запроса"})
	}

	// Валидация
	if req.Email == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Email обязателен"})
	}
	if req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Пароль обязателен"})
	}
	if len(req.Password) < 6 {
		return c.Status(400).JSON(fiber.Map{"error": "Пароль должен содержать минимум 6 символов"})
	}

	user, token, err := h.authService.Register(req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(AuthResponse{
		User:  user,
		Token: token,
	})
}

// Login авторизует пользователя
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверное тело запроса"})
	}

	// Валидация
	if req.Email == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Email обязателен"})
	}
	if req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Пароль обязателен"})
	}

	user, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(AuthResponse{
		User:  user,
		Token: token,
	})
}

// GetProfile возвращает профиль текущего пользователя
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	user, err := h.authService.GetUserProfile(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}

// UpdatePassword обновляет пароль пользователя
func (h *AuthHandler) UpdatePassword(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var req UpdatePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверное тело запроса"})
	}

	// Валидация
	if req.OldPassword == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Текущий пароль обязателен"})
	}
	if req.NewPassword == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Новый пароль обязателен"})
	}
	if len(req.NewPassword) < 6 {
		return c.Status(400).JSON(fiber.Map{"error": "Новый пароль должен содержать минимум 6 символов"})
	}

	err := h.authService.UpdatePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Пароль успешно обновлен"})
}

// UpdateProfile обновляет профиль пользователя
func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var req UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверное тело запроса"})
	}

	user, err := h.authService.UpdateProfile(userID, req.FirstName, req.LastName)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}


