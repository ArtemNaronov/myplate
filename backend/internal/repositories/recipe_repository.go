package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/myplate/backend/internal/models"
	"github.com/myplate/backend/pkg/database"
)

type RecipeRepository struct{}

func NewRecipeRepository() *RecipeRepository {
	return &RecipeRepository{}
}

func (r *RecipeRepository) GetAll() ([]models.Recipe, error) {
	query := `SELECT id, name, description, calories, proteins, fats, carbs, price, cooking_time, servings,
	         meal_type, diet_type, allergens, ingredients, instructions, image_url, created_at, updated_at
	         FROM recipes ORDER BY name`
	
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		var dietType, allergens, instructions []string
		var ingredientsJSON []byte
		var description, mealType, imageURL sql.NullString
		
		var price float64 // Временная переменная для сканирования (поле в БД есть, но не используем)
		err := rows.Scan(
			&recipe.ID, &recipe.Name, &description, &recipe.Calories, &recipe.Proteins,
			&recipe.Fats, &recipe.Carbs, &price, &recipe.CookingTime, &recipe.Servings,
			&mealType, pq.Array(&dietType), pq.Array(&allergens), &ingredientsJSON, pq.Array(&instructions),
			&imageURL, &recipe.CreatedAt, &recipe.UpdatedAt,
		)
		_ = price // Игнорируем цену
		if err != nil {
			return nil, err
		}
		
		recipe.Description = description.String
		recipe.MealType = mealType.String
		recipe.ImageURL = imageURL.String
		recipe.DietType = dietType
		recipe.Allergens = allergens
		recipe.Instructions = instructions
		if ingredientsJSON != nil && len(ingredientsJSON) > 0 {
			json.Unmarshal(ingredientsJSON, &recipe.Ingredients)
		} else {
			recipe.Ingredients = models.Ingredients{}
		}
		
		recipes = append(recipes, recipe)
	}
	
	return recipes, rows.Err()
}

func (r *RecipeRepository) GetByID(id int) (*models.Recipe, error) {
	query := `SELECT id, name, description, calories, proteins, fats, carbs, price, cooking_time, servings,
	         meal_type, diet_type, allergens, ingredients, instructions, image_url, created_at, updated_at
	         FROM recipes WHERE id = $1`
	
	var recipe models.Recipe
	var dietType, allergens, instructions []string
	var ingredientsJSON []byte
	var description, mealType, imageURL sql.NullString
	
	var price float64 // Временная переменная для сканирования (поле в БД есть, но не используем)
	err := database.DB.QueryRow(query, id).Scan(
		&recipe.ID, &recipe.Name, &description, &recipe.Calories, &recipe.Proteins,
		&recipe.Fats, &recipe.Carbs, &price, &recipe.CookingTime, &recipe.Servings,
		&mealType, pq.Array(&dietType), pq.Array(&allergens), &ingredientsJSON, pq.Array(&instructions),
		&imageURL, &recipe.CreatedAt, &recipe.UpdatedAt,
	)
	_ = price // Игнорируем цену
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	recipe.Description = description.String
	recipe.MealType = mealType.String
	recipe.ImageURL = imageURL.String
	recipe.DietType = dietType
	recipe.Allergens = allergens
	recipe.Instructions = instructions
	if ingredientsJSON != nil && len(ingredientsJSON) > 0 {
		json.Unmarshal(ingredientsJSON, &recipe.Ingredients)
	} else {
		recipe.Ingredients = models.Ingredients{}
	}
	
	return &recipe, nil
}

func (r *RecipeRepository) GetFiltered(dietType string, allergies []string, mealTypes []string, maxCalories, maxPrice, maxTime *int) ([]models.Recipe, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if dietType != "" {
		conditions = append(conditions, fmt.Sprintf("$%d = ANY(diet_type)", argIndex))
		args = append(args, dietType)
		argIndex++
	}

	if len(allergies) > 0 {
		var allergyConditions []string
		for _, allergy := range allergies {
			allergyConditions = append(allergyConditions, fmt.Sprintf("$%d != ALL(allergens)", argIndex))
			args = append(args, allergy)
			argIndex++
		}
		conditions = append(conditions, "("+strings.Join(allergyConditions, " AND ")+")")
	}

	if len(mealTypes) > 0 {
		conditions = append(conditions, fmt.Sprintf("meal_type = ANY($%d)", argIndex))
		args = append(args, pq.Array(mealTypes))
		argIndex++
	}

	if maxCalories != nil {
		conditions = append(conditions, fmt.Sprintf("calories <= $%d", argIndex))
		args = append(args, *maxCalories)
		argIndex++
	}

	// maxPrice игнорируется - цены больше не используются

	if maxTime != nil {
		conditions = append(conditions, fmt.Sprintf("cooking_time <= $%d", argIndex))
		args = append(args, *maxTime)
		argIndex++
	}

	query := `SELECT id, name, description, calories, proteins, fats, carbs, price, cooking_time, servings,
	         meal_type, diet_type, allergens, ingredients, instructions, image_url, created_at, updated_at
	         FROM recipes`
	
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	
	query += " ORDER BY name"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		var dietTypeArr, allergensArr, instructions []string
		var ingredientsJSON []byte
		var description, mealType, imageURL sql.NullString
		
		var price float64 // Временная переменная для сканирования (поле в БД есть, но не используем)
		err := rows.Scan(
			&recipe.ID, &recipe.Name, &description, &recipe.Calories, &recipe.Proteins,
			&recipe.Fats, &recipe.Carbs, &price, &recipe.CookingTime, &recipe.Servings,
			&mealType, pq.Array(&dietTypeArr), pq.Array(&allergensArr), &ingredientsJSON, pq.Array(&instructions),
			&imageURL, &recipe.CreatedAt, &recipe.UpdatedAt,
		)
		_ = price // Игнорируем цену
		if err != nil {
			return nil, err
		}
		
		recipe.Description = description.String
		recipe.MealType = mealType.String
		recipe.ImageURL = imageURL.String
		recipe.DietType = dietTypeArr
		recipe.Allergens = allergensArr
		recipe.Instructions = instructions
		if ingredientsJSON != nil && len(ingredientsJSON) > 0 {
			json.Unmarshal(ingredientsJSON, &recipe.Ingredients)
		} else {
			recipe.Ingredients = models.Ingredients{}
		}
		
		recipes = append(recipes, recipe)
	}
	
	return recipes, rows.Err()
}

