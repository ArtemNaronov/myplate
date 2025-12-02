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
		RETURNING id, telegram_id, email, username, first_name, last_name, created_at, updated_at
	`
	
	user := &models.User{}
	var telegramIDNull sql.NullInt64
	var emailNull, usernameNull, firstNameNull, lastNameNull sql.NullString
	err := database.DB.QueryRow(query, telegramID, username, firstName, lastName).Scan(
		&user.ID, &telegramIDNull, &emailNull, &usernameNull, &firstNameNull, &lastNameNull, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if telegramIDNull.Valid {
		telegramIDVal := telegramIDNull.Int64
		user.TelegramID = &telegramIDVal
	}
	user.Email = emailNull.String
	user.Username = usernameNull.String
	user.FirstName = firstNameNull.String
	user.LastName = lastNameNull.String
	return user, nil
}

func (r *UserRepository) GetByTelegramID(telegramID int64) (*models.User, error) {
	query := `SELECT id, telegram_id, email, username, first_name, last_name, created_at, updated_at FROM users WHERE telegram_id = $1`
	
	user := &models.User{}
	var telegramIDNull sql.NullInt64
	var emailNull, usernameNull, firstNameNull, lastNameNull sql.NullString
	err := database.DB.QueryRow(query, telegramID).Scan(
		&user.ID, &telegramIDNull, &emailNull, &usernameNull, &firstNameNull, &lastNameNull, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if telegramIDNull.Valid {
		telegramIDVal := telegramIDNull.Int64
		user.TelegramID = &telegramIDVal
	}
	user.Email = emailNull.String
	user.Username = usernameNull.String
	user.FirstName = firstNameNull.String
	user.LastName = lastNameNull.String
	return user, nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	query := `SELECT id, telegram_id, email, username, first_name, last_name, created_at, updated_at FROM users WHERE id = $1`
	
	user := &models.User{}
	var telegramIDNull sql.NullInt64
	var emailNull, usernameNull, firstNameNull, lastNameNull sql.NullString
	err := database.DB.QueryRow(query, id).Scan(
		&user.ID, &telegramIDNull, &emailNull, &usernameNull, &firstNameNull, &lastNameNull, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if telegramIDNull.Valid {
		telegramIDVal := telegramIDNull.Int64
		user.TelegramID = &telegramIDVal
	}
	user.Email = emailNull.String
	user.Username = usernameNull.String
	user.FirstName = firstNameNull.String
	user.LastName = lastNameNull.String
	return user, nil
}

// GetByEmail получает пользователя по email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `SELECT id, telegram_id, email, username, first_name, last_name, password_hash, created_at, updated_at FROM users WHERE email = $1`
	
	user := &models.User{}
	var telegramIDNull sql.NullInt64
	var emailNull, usernameNull, firstNameNull, lastNameNull, passwordHashNull sql.NullString
	err := database.DB.QueryRow(query, email).Scan(
		&user.ID, &telegramIDNull, &emailNull, &usernameNull, &firstNameNull, &lastNameNull, &passwordHashNull, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if telegramIDNull.Valid {
		telegramIDVal := telegramIDNull.Int64
		user.TelegramID = &telegramIDVal
	}
	user.Email = emailNull.String
	user.Username = usernameNull.String
	user.FirstName = firstNameNull.String
	user.LastName = lastNameNull.String
	user.PasswordHash = passwordHashNull.String
	return user, nil
}

// Create создает нового пользователя с email и паролем
func (r *UserRepository) Create(email, passwordHash, firstName, lastName string) (*models.User, error) {
	query := `
		INSERT INTO users (email, password_hash, first_name, last_name, username)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, telegram_id, email, username, first_name, last_name, created_at, updated_at
	`
	
	user := &models.User{}
	var telegramIDNull sql.NullInt64
	var emailNull, usernameNull, firstNameNull, lastNameNull sql.NullString
	username := email // Используем email как username по умолчанию
	err := database.DB.QueryRow(query, email, passwordHash, firstName, lastName, username).Scan(
		&user.ID, &telegramIDNull, &emailNull, &usernameNull, &firstNameNull, &lastNameNull, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if telegramIDNull.Valid {
		telegramIDVal := telegramIDNull.Int64
		user.TelegramID = &telegramIDVal
	}
	user.Email = emailNull.String
	user.Username = usernameNull.String
	user.FirstName = firstNameNull.String
	user.LastName = lastNameNull.String
	return user, nil
}

// UpdatePassword обновляет пароль пользователя
func (r *UserRepository) UpdatePassword(userID int, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := database.DB.Exec(query, passwordHash, userID)
	return err
}

// UpdateProfile обновляет профиль пользователя
func (r *UserRepository) UpdateProfile(userID int, firstName, lastName string) error {
	query := `UPDATE users SET first_name = $1, last_name = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`
	_, err := database.DB.Exec(query, firstName, lastName, userID)
	return err
}


