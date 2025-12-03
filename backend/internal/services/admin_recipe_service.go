package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/myplate/backend/internal/models"
	"github.com/myplate/backend/internal/repositories"
	"github.com/myplate/backend/pkg/database"
)

type AdminRecipeService struct {
	recipeRepo *repositories.RecipeRepository
}

func NewAdminRecipeService(recipeRepo *repositories.RecipeRepository) *AdminRecipeService {
	return &AdminRecipeService{
		recipeRepo: recipeRepo,
	}
}

// CreateRecipe создает новый рецепт из DTO
func (s *AdminRecipeService) CreateRecipe(dto *models.RecipeImportDTO) (*models.Recipe, error) {
	// Проверяем дубликаты
	ctx := context.Background()
	tx, err := database.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("ошибка при начале транзакции: %w", err)
	}
	defer tx.Rollback()
	
	exists, err := s.recipeRepo.ExistsByName(ctx, tx, dto.Title)
	if err != nil {
		return nil, fmt.Errorf("ошибка при проверке дубликата: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("рецепт с названием '%s' уже существует", dto.Title)
	}
	
	// Преобразуем DTO в модель Recipe
	recipe := s.dtoToRecipe(dto)
	
	// Создаем рецепт в транзакции
	err = s.recipeRepo.CreateInTx(ctx, tx, recipe)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании рецепта: %w", err)
	}
	
	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("ошибка при коммите транзакции: %w", err)
	}
	
	return recipe, nil
}

// ImportRecipes импортирует несколько рецептов
type ImportResult struct {
	Imported int      `json:"imported"`
	Failed   int      `json:"failed"`
	Errors   []string `json:"errors"`
}

func (s *AdminRecipeService) ImportRecipes(recipes []models.RecipeImportDTO) (*ImportResult, error) {
	result := &ImportResult{
		Errors: []string{},
	}
	
	if len(recipes) == 0 {
		return result, nil
	}
	
	// Начинаем транзакцию
	ctx := context.Background()
	tx, err := database.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("ошибка при начале транзакции: %w", err)
	}
	defer tx.Rollback()
	
	// Импортируем рецепты
	for i, dto := range recipes {
		// Проверяем дубликаты по названию
		exists, err := s.recipeRepo.ExistsByName(ctx, tx, dto.Title)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Рецепт %d (%s): ошибка проверки дубликата: %v", i+1, dto.Title, err))
			continue
		}
		
		if exists {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Рецепт %d (%s): рецепт с таким названием уже существует", i+1, dto.Title))
			continue
		}
		
		recipe := s.dtoToRecipe(&dto)
		
		// Создаем рецепт в транзакции
		err = s.recipeRepo.CreateInTx(ctx, tx, recipe)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Рецепт %d (%s): %v", i+1, dto.Title, err))
			continue
		}
		
		result.Imported++
	}
	
	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("ошибка при коммите транзакции: %w", err)
	}
	
	return result, nil
}

// ExportRecipes экспортирует все рецепты
func (s *AdminRecipeService) ExportRecipes() (*models.RecipeExportResponse, error) {
	recipes, err := s.recipeRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении рецептов: %w", err)
	}
	
	exportRecipes := make([]models.RecipeExportDTO, 0, len(recipes))
	for _, recipe := range recipes {
		exportDTO := s.recipeToDTO(&recipe)
		exportRecipes = append(exportRecipes, exportDTO)
	}
	
	return &models.RecipeExportResponse{
		Recipes: exportRecipes,
	}, nil
}

// dtoToRecipe преобразует RecipeImportDTO в Recipe
func (s *AdminRecipeService) dtoToRecipe(dto *models.RecipeImportDTO) *models.Recipe {
	// Преобразуем ингредиенты
	ingredients := make(models.Ingredients, 0, len(dto.Ingredients))
	for _, ing := range dto.Ingredients {
		ingredients = append(ingredients, models.Ingredient{
			Name:     ing.Name,
			Quantity: ing.Amount,
			Unit:     ing.Unit,
		})
	}
	
	// Извлекаем meal_type, diet_type, allergens из tags
	mealType := ""
	dietType := []string{}
	allergens := []string{}
	
	for _, tag := range dto.Tags {
		tagLower := strings.ToLower(tag)
		switch tagLower {
		case "breakfast", "lunch", "dinner", "snack":
			mealType = tagLower
		case "vegetarian", "vegan", "gluten-free":
			dietType = append(dietType, tagLower)
		case "eggs", "dairy", "nuts", "fish", "gluten":
			allergens = append(allergens, tagLower)
		}
	}
	
	if mealType == "" {
		mealType = "lunch" // По умолчанию
	}
	
	servings := dto.Servings
	if servings == 0 {
		servings = 1
	}
	
	cookingTime := dto.CookingTime
	if cookingTime == 0 {
		cookingTime = 30 // По умолчанию
	}
	
	return &models.Recipe{
		Name:         dto.Title,
		Description:  dto.Description,
		Calories:     dto.Calories,
		Proteins:     dto.Proteins,
		Fats:         dto.Fats,
		Carbs:        dto.Carbs,
		CookingTime:  cookingTime,
		Servings:     servings,
		MealType:     mealType,
		DietType:     dietType,
		Allergens:    allergens,
		Ingredients:  ingredients,
		Instructions: dto.Instructions,
	}
}

// recipeToDTO преобразует Recipe в RecipeExportDTO
func (s *AdminRecipeService) recipeToDTO(recipe *models.Recipe) models.RecipeExportDTO {
	// Преобразуем ингредиенты
	ingredients := make([]models.IngredientImport, 0, len(recipe.Ingredients))
	for _, ing := range recipe.Ingredients {
		ingredients = append(ingredients, models.IngredientImport{
			Name:   ing.Name,
			Amount: ing.Quantity,
			Unit:   ing.Unit,
		})
	}
	
	// Формируем tags из meal_type, diet_type, allergens
	tags := []string{}
	if recipe.MealType != "" {
		tags = append(tags, recipe.MealType)
	}
	tags = append(tags, recipe.DietType...)
	tags = append(tags, recipe.Allergens...)
	
	return models.RecipeExportDTO{
		Title:        recipe.Name,
		Description:  recipe.Description,
		Tags:         tags,
		Ingredients:  ingredients,
		Calories:     recipe.Calories,
		Proteins:     recipe.Proteins,
		Fats:         recipe.Fats,
		Carbs:        recipe.Carbs,
		CookingTime:  recipe.CookingTime,
		Servings:     recipe.Servings,
		Instructions: recipe.Instructions,
	}
}

