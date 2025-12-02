package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Recipe struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Calories     int       `json:"calories"`
	Proteins     float64   `json:"proteins"`
	Fats         float64   `json:"fats"`
	Carbs        float64   `json:"carbs"`
	CookingTime  int       `json:"cooking_time"`
	Servings     int       `json:"servings"`
	MealType     string    `json:"meal_type"`
	DietType     []string  `json:"diet_type"`
	Allergens    []string  `json:"allergens"`
	Ingredients  Ingredients `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	ImageURL     string    `json:"image_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Ingredient struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

type Ingredients []Ingredient

func (i *Ingredients) Scan(value interface{}) error {
	if value == nil {
		*i = Ingredients{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(value.(string)), i)
	}
	return json.Unmarshal(bytes, i)
}

func (i Ingredients) Value() (driver.Value, error) {
	if len(i) == 0 {
		return "[]", nil
	}
	return json.Marshal(i)
}


