package models

// RecipeImportDTO - DTO для импорта рецепта из JSON
type RecipeImportDTO struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Tags        []string            `json:"tags"` // diet_type, allergens, meal_type
	Ingredients []IngredientImport  `json:"ingredients"`
	Calories    int                 `json:"calories"`
	Proteins    float64             `json:"proteins"`
	Fats        float64             `json:"fats"`
	Carbs       float64             `json:"carbs"`
	CookingTime int                 `json:"cooking_time,omitempty"`
	Servings    int                 `json:"servings,omitempty"`
	Instructions []string            `json:"instructions,omitempty"`
}

type IngredientImport struct {
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
	Unit     string  `json:"unit"`
}

// RecipeExportDTO - DTO для экспорта рецепта в JSON
type RecipeExportDTO struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Tags        []string            `json:"tags"`
	Ingredients []IngredientImport  `json:"ingredients"`
	Calories    int                 `json:"calories"`
	Proteins    float64             `json:"proteins"`
	Fats        float64             `json:"fats"`
	Carbs       float64             `json:"carbs"`
	CookingTime int                 `json:"cooking_time,omitempty"`
	Servings    int                 `json:"servings,omitempty"`
	Instructions []string            `json:"instructions,omitempty"`
}

// RecipeImportRequest - запрос на импорт рецептов
type RecipeImportRequest struct {
	Recipes []RecipeImportDTO `json:"recipes"`
}

// RecipeExportResponse - ответ на экспорт рецептов
type RecipeExportResponse struct {
	Recipes []RecipeExportDTO `json:"recipes"`
}

