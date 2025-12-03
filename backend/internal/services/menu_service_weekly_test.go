package services

import (
	"testing"

	"github.com/myplate/backend/internal/models"
)

func TestMenuService_SelectRecipeForMeal(t *testing.T) {
	service := &MenuService{}
	
	// Создаем тестовые рецепты
	scoredRecipes := []ScoredRecipe{
		{
			Recipe: models.Recipe{
				ID:         1,
				Name:       "Рецепт 1",
				Calories:   500,
				CookingTime: 20,
			},
			Score: 0.8,
		},
		{
			Recipe: models.Recipe{
				ID:         2,
				Name:       "Рецепт 2",
				Calories:   600,
				CookingTime: 25,
			},
			Score: 0.9,
		},
		{
			Recipe: models.Recipe{
				ID:         3,
				Name:       "Рецепт 3",
				Calories:   550,
				CookingTime: 22,
			},
			Score: 0.7,
		},
	}
	
	targetCalories := 550
	excludedIDs := make(map[int]bool)
	
	selected := service.selectRecipeForMeal(scoredRecipes, targetCalories, excludedIDs)
	
	if selected == nil {
		t.Fatal("Не выбран рецепт")
	}
	
	// Должен быть выбран рецепт 3 (ближайший к целевым калориям)
	if selected.ID != 3 {
		t.Errorf("Ожидался рецепт с ID 3, получен ID %d", selected.ID)
	}
	
	// Тест с исключенными рецептами
	excludedIDs[3] = true
	selected = service.selectRecipeForMeal(scoredRecipes, targetCalories, excludedIDs)
	
	if selected == nil {
		t.Fatal("Не выбран рецепт")
	}
	
	// Должен быть выбран рецепт 1 или 2 (3 исключен)
	if selected.ID == 3 {
		t.Errorf("Рецепт с ID 3 должен быть исключен")
	}
}

func TestMenuService_RecipeToDTO(t *testing.T) {
	service := &MenuService{}
	
	recipe := &models.Recipe{
		ID:           1,
		Name:         "Тестовый рецепт",
		Description:  "Описание",
		Calories:     500,
		Proteins:     20.0,
		Fats:         15.0,
		Carbs:        60.0,
		CookingTime:  30,
		Servings:     2,
		MealType:     "breakfast",
		Ingredients:  models.Ingredients{},
		Instructions: []string{"Шаг 1", "Шаг 2"},
	}
	
	dto := service.recipeToDTO(recipe)
	
	if dto == nil {
		t.Fatal("DTO не создан")
	}
	
	if dto.ID != recipe.ID {
		t.Errorf("Ожидался ID %d, получен %d", recipe.ID, dto.ID)
	}
	
	if dto.Name != recipe.Name {
		t.Errorf("Ожидалось название %s, получено %s", recipe.Name, dto.Name)
	}
	
	if dto.Calories != recipe.Calories {
		t.Errorf("Ожидалось калорий %d, получено %d", recipe.Calories, dto.Calories)
	}
}

