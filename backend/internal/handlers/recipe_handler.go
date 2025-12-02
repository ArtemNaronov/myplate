package handlers

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/myplate/backend/internal/services"
)

type RecipeHandler struct {
	recipeService *services.RecipeService
}

func NewRecipeHandler(recipeService *services.RecipeService) *RecipeHandler {
	return &RecipeHandler{
		recipeService: recipeService,
	}
}

func (h *RecipeHandler) GetAll(c *fiber.Ctx) error {
	recipes, err := h.recipeService.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(recipes)
}

func (h *RecipeHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный ID рецепта"})
	}
	
	recipe, err := h.recipeService.GetByID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
	if recipe == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Рецепт не найден"})
	}
	
	return c.JSON(recipe)
}


