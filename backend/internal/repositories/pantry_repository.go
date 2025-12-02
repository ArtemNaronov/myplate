package repositories

import (
	"database/sql"
	"github.com/myplate/backend/internal/models"
	"github.com/myplate/backend/pkg/database"
)

type PantryRepository struct{}

func NewPantryRepository() *PantryRepository {
	return &PantryRepository{}
}

func (r *PantryRepository) GetByUserID(userID int) ([]models.PantryItem, error) {
	query := `SELECT id, user_id, name, quantity, unit, created_at, updated_at 
	         FROM pantry_items WHERE user_id = $1 ORDER BY name`
	
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.PantryItem
	for rows.Next() {
		var item models.PantryItem
		err := rows.Scan(
			&item.ID, &item.UserID, &item.Name, &item.Quantity, &item.Unit,
			&item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	
	return items, rows.Err()
}

func (r *PantryRepository) Create(item *models.PantryItem) error {
	query := `INSERT INTO pantry_items (user_id, name, quantity, unit)
	         VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	
	err := database.DB.QueryRow(query, item.UserID, item.Name, item.Quantity, item.Unit).Scan(
		&item.ID, &item.CreatedAt, &item.UpdatedAt,
	)
	return err
}

func (r *PantryRepository) Delete(id, userID int) error {
	query := `DELETE FROM pantry_items WHERE id = $1 AND user_id = $2`
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


