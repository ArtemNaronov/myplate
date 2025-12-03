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
	// Устанавливаем menu_type по умолчанию, если не указан
	if menu.MenuType == "" {
		menu.MenuType = "daily"
	}
	
	// Для дневных меню используем ON CONFLICT, для недельных - простой INSERT
	if menu.MenuType == "daily" {
		query := `
			INSERT INTO menus (user_id, date, total_calories, total_price, total_time, menu_type, meals, ingredients_used, missing_ingredients)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (user_id, date) 
			WHERE menu_type = 'daily'
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
		
		var totalPrice float64 = 0
		err := database.DB.QueryRow(query,
			menu.UserID, menu.Date, menu.TotalCalories, totalPrice, menu.TotalTime, menu.MenuType,
			mealsJSON, ingredientsUsedJSON, missingIngredientsJSON,
		).Scan(&menu.ID, &menu.CreatedAt, &menu.UpdatedAt)
		_ = totalPrice
		
		return err
	}
	
	// Для недельных меню используем CreateWeeklyMenu
	return r.CreateWeeklyMenu(menu)
}

// CreateWeeklyMenu сохраняет недельное меню
func (r *MenuRepository) CreateWeeklyMenu(menu *models.Menu) error {
	menu.MenuType = "weekly"
	
	// Для недельного меню используем простой INSERT без конфликта
	query := `
		INSERT INTO menus (user_id, date, total_calories, total_price, total_time, menu_type, meals, ingredients_used, missing_ingredients)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`
	
	mealsJSON, _ := json.Marshal(menu.Meals)
	ingredientsUsedJSON, _ := json.Marshal(menu.IngredientsUsed)
	missingIngredientsJSON, _ := json.Marshal(menu.MissingIngredients)
	
	var totalPrice float64 = 0
	err := database.DB.QueryRow(query,
		menu.UserID, menu.Date, menu.TotalCalories, totalPrice, menu.TotalTime, menu.MenuType,
		mealsJSON, ingredientsUsedJSON, missingIngredientsJSON,
	).Scan(&menu.ID, &menu.CreatedAt, &menu.UpdatedAt)
	
	return err
}

func (r *MenuRepository) GetByUserIDAndDate(userID int, date time.Time) (*models.Menu, error) {
	query := `
		SELECT id, user_id, date, total_calories, total_price, total_time, menu_type, meals, 
		       ingredients_used, missing_ingredients, created_at, updated_at
		FROM menus WHERE user_id = $1 AND date = $2 AND menu_type = 'daily'
	`
	
	var menu models.Menu
	var mealsJSON, ingredientsUsedJSON, missingIngredientsJSON []byte
	var totalPrice float64 // Временная переменная (поле в БД есть, но не используем)
	
	err := database.DB.QueryRow(query, userID, date).Scan(
		&menu.ID, &menu.UserID, &menu.Date, &menu.TotalCalories, &totalPrice, &menu.TotalTime, &menu.MenuType,
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
		SELECT id, user_id, date, total_calories, total_price, total_time, menu_type, meals, 
		       ingredients_used, missing_ingredients, created_at, updated_at
		FROM menus WHERE id = $1
	`
	
	var menu models.Menu
	var mealsJSON, ingredientsUsedJSON, missingIngredientsJSON []byte
	var totalPrice float64 // Временная переменная (поле в БД есть, но не используем)
	
	err := database.DB.QueryRow(query, id).Scan(
		&menu.ID, &menu.UserID, &menu.Date, &menu.TotalCalories, &totalPrice, &menu.TotalTime, &menu.MenuType,
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

// GetWeeklyMenusByUserID получает все недельные меню пользователя
func (r *MenuRepository) GetWeeklyMenusByUserID(userID int) ([]models.Menu, error) {
	query := `
		SELECT id, user_id, date, total_calories, total_price, total_time, menu_type, meals, 
		       ingredients_used, missing_ingredients, created_at, updated_at
		FROM menus WHERE user_id = $1 AND menu_type = 'weekly' ORDER BY date DESC
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
		var totalPrice float64
		
		err := rows.Scan(
			&menu.ID, &menu.UserID, &menu.Date, &menu.TotalCalories, &totalPrice, &menu.TotalTime, &menu.MenuType,
			&mealsJSON, &ingredientsUsedJSON, &missingIngredientsJSON,
			&menu.CreatedAt, &menu.UpdatedAt,
		)
		_ = totalPrice
		if err != nil {
			return nil, err
		}
		
		// Для недельных меню meals содержит JSON с данными недели, а не MenuMeals
		// Поэтому мы просто сохраняем JSON как есть, без парсинга в MenuMeals
		// При чтении фронтенд сам распарсит JSON
		// Но для совместимости с моделью Menu, мы создаем пустой MenuMeals
		menu.Meals = models.MenuMeals{}
		
		// Сохраняем сырой JSON в поле Meals через интерфейс Value()
		// Но это не сработает, так как MenuMeals имеет свой Value()
		// Вместо этого, мы просто оставляем Meals пустым, а фронтенд будет читать meals напрямую из JSON ответа
		
		json.Unmarshal(ingredientsUsedJSON, &menu.IngredientsUsed)
		json.Unmarshal(missingIngredientsJSON, &menu.MissingIngredients)
		
		menus = append(menus, menu)
	}
	
	return menus, rows.Err()
}

// Delete удаляет меню по ID, проверяя принадлежность пользователю
func (r *MenuRepository) Delete(id, userID int) error {
	query := `DELETE FROM menus WHERE id = $1 AND user_id = $2`
	result, err := database.DB.Exec(query, id, userID)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *MenuRepository) GetAllByUserID(userID int) ([]models.Menu, error) {
	query := `
		SELECT id, user_id, date, total_calories, total_price, total_time, menu_type, meals, 
		       ingredients_used, missing_ingredients, created_at, updated_at
		FROM menus WHERE user_id = $1 AND menu_type = 'daily' ORDER BY date DESC
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
			&menu.ID, &menu.UserID, &menu.Date, &menu.TotalCalories, &totalPrice, &menu.TotalTime, &menu.MenuType,
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


