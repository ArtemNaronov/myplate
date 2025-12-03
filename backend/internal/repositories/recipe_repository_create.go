package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"
	"github.com/myplate/backend/internal/models"
	"github.com/myplate/backend/pkg/database"
)

// Create создает новый рецепт
func (r *RecipeRepository) Create(recipe *models.Recipe) (*models.Recipe, error) {
	ctx := context.Background()
	tx, err := database.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("ошибка при начале транзакции: %w", err)
	}
	defer tx.Rollback()
	
	err = r.CreateInTx(ctx, tx, recipe)
	if err != nil {
		return nil, err
	}
	
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("ошибка при коммите транзакции: %w", err)
	}
	
	return recipe, nil
}

// CreateInTx создает рецепт в транзакции
func (r *RecipeRepository) CreateInTx(ctx context.Context, tx *sql.Tx, recipe *models.Recipe) error {
	query := `
		INSERT INTO recipes (name, description, calories, proteins, fats, carbs, cooking_time, servings,
		                     meal_type, diet_type, allergens, ingredients, instructions, image_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at, updated_at
	`
	
	ingredientsJSON, err := json.Marshal(recipe.Ingredients)
	if err != nil {
		return fmt.Errorf("ошибка при сериализации ингредиентов: %w", err)
	}
	
	var description, mealType, imageURL sql.NullString
	if recipe.Description != "" {
		description.String = recipe.Description
		description.Valid = true
	}
	if recipe.MealType != "" {
		mealType.String = recipe.MealType
		mealType.Valid = true
	}
	if recipe.ImageURL != "" {
		imageURL.String = recipe.ImageURL
		imageURL.Valid = true
	}
	
	var price float64 = 0 // Цена не используется, но поле есть в БД
	
	err = tx.QueryRowContext(ctx, query,
		recipe.Name, description, recipe.Calories, recipe.Proteins, recipe.Fats, recipe.Carbs,
		recipe.CookingTime, recipe.Servings, mealType, pq.Array(recipe.DietType),
		pq.Array(recipe.Allergens), ingredientsJSON, pq.Array(recipe.Instructions), imageURL,
	).Scan(&recipe.ID, &recipe.CreatedAt, &recipe.UpdatedAt)
	_ = price // Игнорируем цену
	
	if err != nil {
		return fmt.Errorf("ошибка при создании рецепта: %w", err)
	}
	
	return nil
}

// ExistsByName проверяет, существует ли рецепт с таким названием
func (r *RecipeRepository) ExistsByName(ctx context.Context, tx *sql.Tx, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM recipes WHERE LOWER(name) = LOWER($1))`
	
	var exists bool
	err := tx.QueryRowContext(ctx, query, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("ошибка при проверке существования рецепта: %w", err)
	}
	
	return exists, nil
}

