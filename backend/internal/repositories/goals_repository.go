package repositories

import (
	"database/sql"
	"github.com/myplate/backend/internal/models"
	"github.com/myplate/backend/pkg/database"
)

type GoalsRepository struct{}

func NewGoalsRepository() *GoalsRepository {
	return &GoalsRepository{}
}

func (r *GoalsRepository) CreateOrUpdate(goals *models.UserGoals) error {
	query := `
		INSERT INTO user_goals (
			user_id, daily_calories, target_proteins, target_fats, target_carbs,
			protein_ratio, fat_ratio, carb_ratio
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id) 
		DO UPDATE SET 
			daily_calories = EXCLUDED.daily_calories,
			target_proteins = EXCLUDED.target_proteins,
			target_fats = EXCLUDED.target_fats,
			target_carbs = EXCLUDED.target_carbs,
			protein_ratio = EXCLUDED.protein_ratio,
			fat_ratio = EXCLUDED.fat_ratio,
			carb_ratio = EXCLUDED.carb_ratio,
			updated_at = CURRENT_TIMESTAMP
	`
	
	_, err := database.DB.Exec(query,
		goals.UserID, goals.DailyCalories, goals.TargetProteins, goals.TargetFats, goals.TargetCarbs,
		goals.ProteinRatio, goals.FatRatio, goals.CarbRatio,
	)
	return err
}

func (r *GoalsRepository) GetByUserID(userID int) (*models.UserGoals, error) {
	query := `
		SELECT id, user_id, daily_calories, target_proteins, target_fats, target_carbs,
		       protein_ratio, fat_ratio, carb_ratio, created_at, updated_at
		FROM user_goals WHERE user_id = $1
	`
	
	goals := &models.UserGoals{}
	err := database.DB.QueryRow(query, userID).Scan(
		&goals.ID, &goals.UserID, &goals.DailyCalories, &goals.TargetProteins, &goals.TargetFats, &goals.TargetCarbs,
		&goals.ProteinRatio, &goals.FatRatio, &goals.CarbRatio, &goals.CreatedAt, &goals.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return goals, nil
}


