package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/myplate/backend/internal/models"
	"github.com/myplate/backend/internal/services"
)

type AdminRecipeHandler struct {
	adminRecipeService *services.AdminRecipeService
}

func NewAdminRecipeHandler(adminRecipeService *services.AdminRecipeService) *AdminRecipeHandler {
	return &AdminRecipeHandler{
		adminRecipeService: adminRecipeService,
	}
}

// Create создает новый рецепт
func (h *AdminRecipeHandler) Create(c *fiber.Ctx) error {
	var req models.RecipeImportDTO
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверное тело запроса"})
	}
	
	recipe, err := h.adminRecipeService.CreateRecipe(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
	return c.Status(201).JSON(recipe)
}

// Import импортирует рецепты из JSON
func (h *AdminRecipeHandler) Import(c *fiber.Ctx) error {
	var req models.RecipeImportRequest
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверное тело запроса"})
	}
	
	if len(req.Recipes) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Список рецептов пуст"})
	}
	
	result, err := h.adminRecipeService.ImportRecipes(req.Recipes)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
	return c.Status(201).JSON(fiber.Map{
		"imported": result.Imported,
		"failed":   result.Failed,
		"errors":   result.Errors,
	})
}

// Export экспортирует все рецепты в JSON
func (h *AdminRecipeHandler) Export(c *fiber.Ctx) error {
	exportData, err := h.adminRecipeService.ExportRecipes()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
	return c.JSON(exportData)
}

