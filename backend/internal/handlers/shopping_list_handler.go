package handlers

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/myplate/backend/internal/services"
)

type ShoppingListHandler struct {
	shoppingService *services.ShoppingListService
}

func NewShoppingListHandler(shoppingService *services.ShoppingListService) *ShoppingListHandler {
	return &ShoppingListHandler{
		shoppingService: shoppingService,
	}
}

func (h *ShoppingListHandler) GetByMenuID(c *fiber.Ctx) error {
	menuID, err := strconv.Atoi(c.Params("menu_id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID меню"})
	}
	
	list, err := h.shoppingService.GetByMenuID(menuID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
	if list == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Список покупок не найден"})
	}
	
	return c.JSON(list)
}


