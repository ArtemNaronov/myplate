package repositories

import (
	"database/sql"
	"encoding/json"
	"github.com/myplate/backend/internal/models"
	"github.com/myplate/backend/pkg/database"
)

type ShoppingListRepository struct{}

func NewShoppingListRepository() *ShoppingListRepository {
	return &ShoppingListRepository{}
}

func (r *ShoppingListRepository) CreateOrUpdate(list *models.ShoppingList) error {
	// Try to get existing
	existing, err := r.GetByMenuID(list.MenuID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	
	itemsJSON, _ := json.Marshal(list.Items)
	
	if existing != nil {
		// Update existing
		updateQuery := `UPDATE shopping_lists SET items = $1, updated_at = CURRENT_TIMESTAMP 
		               WHERE id = $2 RETURNING created_at, updated_at`
		return database.DB.QueryRow(updateQuery, itemsJSON, existing.ID).Scan(
			&list.CreatedAt, &list.UpdatedAt,
		)
	}
	
	// Create new
	query := `INSERT INTO shopping_lists (user_id, menu_id, items)
	         VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	
	return database.DB.QueryRow(query, list.UserID, list.MenuID, itemsJSON).Scan(
		&list.ID, &list.CreatedAt, &list.UpdatedAt,
	)
}

func (r *ShoppingListRepository) GetByMenuID(menuID int) (*models.ShoppingList, error) {
	query := `SELECT id, user_id, menu_id, items, created_at, updated_at 
	         FROM shopping_lists WHERE menu_id = $1`
	
	var list models.ShoppingList
	var itemsJSON []byte
	
	err := database.DB.QueryRow(query, menuID).Scan(
		&list.ID, &list.UserID, &list.MenuID, &itemsJSON, &list.CreatedAt, &list.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	json.Unmarshal(itemsJSON, &list.Items)
	return &list, nil
}

