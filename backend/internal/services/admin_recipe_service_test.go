package services

import (
	"testing"

	"github.com/myplate/backend/internal/models"
)

func TestAdminRecipeService_DtoToRecipe(t *testing.T) {
	service := &AdminRecipeService{}
	
	dto := &models.RecipeImportDTO{
		Title:       "Тестовый рецепт",
		Description: "Описание",
		Tags:        []string{"breakfast", "vegetarian", "eggs"},
		Ingredients: []models.IngredientImport{
			{Name: "Яйца", Amount: 2, Unit: "шт"},
			{Name: "Молоко", Amount: 100, Unit: "мл"},
		},
		Calories:    350,
		Proteins:    20.0,
		Fats:        18.0,
		Carbs:       28.0,
		CookingTime: 10,
		Servings:    1,
		Instructions: []string{"Шаг 1", "Шаг 2"},
	}
	
	recipe := service.dtoToRecipe(dto)
	
	if recipe == nil {
		t.Fatal("Рецепт не создан")
	}
	
	if recipe.Name != dto.Title {
		t.Errorf("Ожидалось название %s, получено %s", dto.Title, recipe.Name)
	}
	
	if recipe.Calories != dto.Calories {
		t.Errorf("Ожидалось калорий %d, получено %d", dto.Calories, recipe.Calories)
	}
	
	if recipe.MealType != "breakfast" {
		t.Errorf("Ожидался meal_type 'breakfast', получен '%s'", recipe.MealType)
	}
	
	// Проверяем, что теги правильно преобразованы
	if len(recipe.DietType) == 0 || recipe.DietType[0] != "vegetarian" {
		t.Errorf("Ожидался diet_type 'vegetarian', получен %v", recipe.DietType)
	}
	
	if len(recipe.Allergens) == 0 || recipe.Allergens[0] != "eggs" {
		t.Errorf("Ожидался аллерген 'eggs', получен %v", recipe.Allergens)
	}
	
	if len(recipe.Ingredients) != 2 {
		t.Errorf("Ожидалось 2 ингредиента, получено %d", len(recipe.Ingredients))
	}
}

func TestAdminRecipeService_RecipeToDTO(t *testing.T) {
	service := &AdminRecipeService{}
	
	recipe := &models.Recipe{
		ID:           1,
		Name:         "Тестовый рецепт",
		Description:  "Описание",
		Calories:     500,
		Proteins:     25.0,
		Fats:         20.0,
		Carbs:        50.0,
		CookingTime:  30,
		Servings:     2,
		MealType:     "lunch",
		DietType:     []string{"vegetarian"},
		Allergens:    []string{"dairy"},
		Ingredients: models.Ingredients{
			{Name: "Помидоры", Quantity: 200, Unit: "г"},
		},
		Instructions: []string{"Шаг 1"},
	}
	
	dto := service.recipeToDTO(recipe)
	
	if dto.Title != recipe.Name {
		t.Errorf("Ожидалось название %s, получено %s", recipe.Name, dto.Title)
	}
	
	// Проверяем, что теги правильно собраны
	if len(dto.Tags) < 3 {
		t.Errorf("Ожидалось минимум 3 тега, получено %d", len(dto.Tags))
	}
	
	// Проверяем наличие meal_type в тегах
	foundMealType := false
	for _, tag := range dto.Tags {
		if tag == "lunch" {
			foundMealType = true
			break
		}
	}
	if !foundMealType {
		t.Errorf("Тег 'lunch' не найден в тегах")
	}
	
	if len(dto.Ingredients) != 1 {
		t.Errorf("Ожидался 1 ингредиент, получено %d", len(dto.Ingredients))
	}
}

