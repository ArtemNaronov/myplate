package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Menu struct {
	ID                 int       `json:"id"`
	UserID             int       `json:"user_id"`
	Date               time.Time `json:"date"`
	TotalCalories      int       `json:"total_calories"`
	TotalTime          int       `json:"total_time"`
	Meals              MenuMeals `json:"meals"`
	IngredientsUsed    Ingredients `json:"ingredients_used,omitempty"`
	MissingIngredients Ingredients `json:"missing_ingredients,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type MenuMeal struct {
	RecipeID  int     `json:"recipe_id"`
	MealType  string  `json:"meal_type"`
	Calories  int     `json:"calories"`
	Time      int     `json:"time"`
}

type MenuMeals []MenuMeal

func (m *MenuMeals) Scan(value interface{}) error {
	if value == nil {
		*m = MenuMeals{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(value.(string)), m)
	}
	return json.Unmarshal(bytes, m)
}

func (m MenuMeals) Value() (driver.Value, error) {
	if len(m) == 0 {
		return "[]", nil
	}
	return json.Marshal(m)
}

type MenuGenerateRequest struct {
	UserID            int     `json:"user_id"`
	TargetCalories    int     `json:"target_calories"`
	DietType          string  `json:"diet_type,omitempty"`
	Allergies         []string `json:"allergies,omitempty"`
	MaxTotalTime      int     `json:"max_total_time,omitempty"`
	MaxTimePerMeal    int     `json:"max_time_per_meal,omitempty"`
	SpeedLevel        string  `json:"speed_level,omitempty"` // fast, normal, slow
	ConsiderPantry    bool    `json:"consider_pantry"`
	PantryImportance  string  `json:"pantry_importance"` // strict, prefer, ignore
}


