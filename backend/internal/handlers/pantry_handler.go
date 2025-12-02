package handlers

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/myplate/backend/internal/models"
	"github.com/myplate/backend/internal/services"
)

type PantryHandler struct {
	pantryService *services.PantryService
}

func NewPantryHandler(pantryService *services.PantryService) *PantryHandler {
	return &PantryHandler{
		pantryService: pantryService,
	}
}

func (h *PantryHandler) GetAll(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	
	items, err := h.pantryService.GetByUserID(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
	return c.JSON(items)
}

func (h *PantryHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	
	var item models.PantryItem
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверное тело запроса"})
	}
	
	if err := h.pantryService.Create(userID, &item); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
	return c.Status(201).JSON(item)
}

func (h *PantryHandler) Delete(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID продукта"})
	}
	
	if err := h.pantryService.Delete(userID, id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
		return c.JSON(fiber.Map{"message": "Продукт успешно удалён"})
}


