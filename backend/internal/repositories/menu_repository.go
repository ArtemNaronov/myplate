package repositories

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/myplate/backend/internal/models"
	"github.com/myplate/backend/pkg/database"
)

type MenuRepository struct{}

func NewMenuRepository() *MenuRepository {
	return &MenuRepository{}
}

func (r *MenuRepository) Create(menu *models.Menu) error {
	query := `
		INSERT INTO menus (user_id, date, total_calories, total_price, total_time, meals, ingredients_used, missing_ingredients)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id, date) 
		DO UPDATE SET 
			total_calories = EXCLUDED.total_calories,
			total_price = EXCLUDED.total_price,
			total_time = EXCLUDED.total_time,
			meals = EXCLUDED.meals,
			ingredients_used = EXCLUDED.ingredients_used,
			missing_ingredients = EXCLUDED.missing_ingredients,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, created_at, updated_at
	`
	
	mealsJSON, _ := json.Marshal(menu.Meals)
	ingredientsUsedJSON, _ := json.Marshal(menu.IngredientsUsed)
	missingIngredientsJSON, _ := json.Marshal(menu.MissingIngredients)
	
	var totalPrice float64 = 0 // Временная переменная (поле в БД есть, но не используем)
	err := database.DB.QueryRow(query,
		menu.UserID, menu.Date, menu.TotalCalories, totalPrice, menu.TotalTime,
		mealsJSON, ingredientsUsedJSON, missingIngredientsJSON,
	).Scan(&menu.ID, &menu.CreatedAt, &menu.UpdatedAt)
	_ = totalPrice // Игнорируем цену
	
	return err
}

func (r *MenuRepository) GetByUserIDAndDate(userID int, date time.Time) (*models.Menu, error) {
	query := `
		SELECT id, user_id, date, total_calories, total_price, total_time, meals, 
		       ingredients_used, missing_ingredients, created_at, updated_at
		FROM menus WHERE user_id = $1 AND date = $2
	`
	
	var menu models.Menu
	var mealsJSON, ingredientsUsedJSON, missingIngredientsJSON []byte
	var totalPrice float64 // Временная переменная (поле в БД есть, но не используем)
	
	err := database.DB.QueryRow(query, userID, date).Scan(
		&menu.ID, &menu.UserID, &menu.Date, &menu.TotalCalories, &totalPrice, &menu.TotalTime,
		&mealsJSON, &ingredientsUsedJSON, &missingIngredientsJSON,
		&menu.CreatedAt, &menu.UpdatedAt,
	)
	_ = totalPrice // Игнорируем цену
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	json.Unmarshal(mealsJSON, &menu.Meals)
	json.Unmarshal(ingredientsUsedJSON, &menu.IngredientsUsed)
	json.Unmarshal(missingIngredientsJSON, &menu.MissingIngredients)
	
	return &menu, nil
}

func (r *MenuRepository) GetByID(id int) (*models.Menu, error) {
	query := `
		SELECT id, user_id, date, total_calories, total_price, total_time, meals, 
		       ingredients_used, missing_ingredients, created_at, updated_at
		FROM menus WHERE id = $1
	`
	
	var menu models.Menu
	var mealsJSON, ingredientsUsedJSON, missingIngredientsJSON []byte
	var totalPrice float64 // Временная переменная (поле в БД есть, но не используем)
	
	err := database.DB.QueryRow(query, id).Scan(
		&menu.ID, &menu.UserID, &menu.Date, &menu.TotalCalories, &totalPrice, &menu.TotalTime,
		&mealsJSON, &ingredientsUsedJSON, &missingIngredientsJSON,
		&menu.CreatedAt, &menu.UpdatedAt,
	)
	_ = totalPrice // Игнорируем цену
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	json.Unmarshal(mealsJSON, &menu.Meals)
	json.Unmarshal(ingredientsUsedJSON, &menu.IngredientsUsed)
	json.Unmarshal(missingIngredientsJSON, &menu.MissingIngredients)
	
	return &menu, nil
}

func (r *MenuRepository) GetAllByUserID(userID int) ([]models.Menu, error) {
	query := `
		SELECT id, user_id, date, total_calories, total_price, total_time, meals, 
		       ingredients_used, missing_ingredients, created_at, updated_at
		FROM menus WHERE user_id = $1 ORDER BY date DESC
	`
	
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var menus []models.Menu
	for rows.Next() {
		var menu models.Menu
		var mealsJSON, ingredientsUsedJSON, missingIngredientsJSON []byte
		var totalPrice float64 // Временная переменная (поле в БД есть, но не используем)
		
		err := rows.Scan(
			&menu.ID, &menu.UserID, &menu.Date, &menu.TotalCalories, &totalPrice, &menu.TotalTime,
			&mealsJSON, &ingredientsUsedJSON, &missingIngredientsJSON,
			&menu.CreatedAt, &menu.UpdatedAt,
		)
		_ = totalPrice // Игнорируем цену
		if err != nil {
			return nil, err
		}
		
		json.Unmarshal(mealsJSON, &menu.Meals)
		json.Unmarshal(ingredientsUsedJSON, &menu.IngredientsUsed)
		json.Unmarshal(missingIngredientsJSON, &menu.MissingIngredients)
		
		menus = append(menus, menu)
	}
	
	return menus, rows.Err()
}


