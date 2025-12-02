package repositories

import (
	"database/sql"
	"github.com/myplate/backend/internal/models"
	"github.com/myplate/backend/pkg/database"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) CreateOrUpdate(telegramID int64, username, firstName, lastName string) (*models.User, error) {
	query := `
		INSERT INTO users (telegram_id, username, first_name, last_name)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (telegram_id) 
		DO UPDATE SET 
			username = EXCLUDED.username,
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, telegram_id, username, first_name, last_name, created_at, updated_at
	`
	
	user := &models.User{}
	var usernameNull, firstNameNull, lastNameNull sql.NullString
	err := database.DB.QueryRow(query, telegramID, username, firstName, lastName).Scan(
		&user.ID, &user.TelegramID, &usernameNull, &firstNameNull, &lastNameNull, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	user.Username = usernameNull.String
	user.FirstName = firstNameNull.String
	user.LastName = lastNameNull.String
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByTelegramID(telegramID int64) (*models.User, error) {
	query := `SELECT id, telegram_id, username, first_name, last_name, created_at, updated_at FROM users WHERE telegram_id = $1`
	
	user := &models.User{}
	var usernameNull, firstNameNull, lastNameNull sql.NullString
	err := database.DB.QueryRow(query, telegramID).Scan(
		&user.ID, &user.TelegramID, &usernameNull, &firstNameNull, &lastNameNull, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user.Username = usernameNull.String
	user.FirstName = firstNameNull.String
	user.LastName = lastNameNull.String
	return user, nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	query := `SELECT id, telegram_id, username, first_name, last_name, created_at, updated_at FROM users WHERE id = $1`
	
	user := &models.User{}
	var usernameNull, firstNameNull, lastNameNull sql.NullString
	err := database.DB.QueryRow(query, id).Scan(
		&user.ID, &user.TelegramID, &usernameNull, &firstNameNull, &lastNameNull, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user.Username = usernameNull.String
	user.FirstName = firstNameNull.String
	user.LastName = lastNameNull.String
	return user, nil
}


