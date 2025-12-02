package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type ShoppingList struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	MenuID    int       `json:"menu_id"`
	Items     ShoppingItems `json:"items"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ShoppingItem struct {
	Name     string   `json:"name"`
	Quantity float64  `json:"quantity"`
	Unit     string   `json:"unit"`
	Reason   []string `json:"reason"` // meal types that need this ingredient
}

type ShoppingItems []ShoppingItem

func (s *ShoppingItems) Scan(value interface{}) error {
	if value == nil {
		*s = ShoppingItems{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(value.(string)), s)
	}
	return json.Unmarshal(bytes, s)
}

func (s ShoppingItems) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	return json.Marshal(s)
}

