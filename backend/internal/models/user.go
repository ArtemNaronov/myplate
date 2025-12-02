package models

import "time"

type User struct {
	ID         int       `json:"id"`
	TelegramID int64     `json:"telegram_id"`
	Username   string    `json:"username"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserGoals struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	DailyCalories int       `json:"daily_calories"`
	TargetProteins float64  `json:"target_proteins"`
	TargetFats    float64  `json:"target_fats"`
	TargetCarbs   float64  `json:"target_carbs"`
	ProteinRatio  float64  `json:"protein_ratio"`
	FatRatio      float64  `json:"fat_ratio"`
	CarbRatio     float64  `json:"carb_ratio"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}


