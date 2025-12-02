package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/myplate/backend/internal/models"
	"github.com/myplate/backend/internal/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) SetGoals(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	
	var goals models.UserGoals
	if err := c.BodyParser(&goals); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверное тело запроса"})
	}
	
	if err := h.userService.SetGoals(userID, &goals); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
		return c.JSON(fiber.Map{"message": "Цели успешно обновлены"})
}

func (h *UserHandler) GetGoals(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	
	goals, err := h.userService.GetGoals(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
	if goals == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Цели не найдены"})
	}
	
	return c.JSON(goals)
}


