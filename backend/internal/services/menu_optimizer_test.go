package services

import (
	"testing"

	"github.com/myplate/backend/internal/models"
)

func TestMenuOptimizer_CalculateWeeklyMacros(t *testing.T) {
	optimizer := NewMenuOptimizer()
	
	// Создаем тестовое недельное меню
	weeklyMenu := &models.WeeklyMenu{
		Week: []models.WeeklyDayMenu{
			{
				Day: 1,
				Breakfast: &models.RecipeDTO{
					Calories: 500,
					Proteins: 20.0,
					Fats:     15.0,
					Carbs:    60.0,
					Servings: 1,
				},
				Lunch: &models.RecipeDTO{
					Calories: 800,
					Proteins: 40.0,
					Fats:     30.0,
					Carbs:    80.0,
					Servings: 1,
				},
				Dinner: &models.RecipeDTO{
					Calories: 700,
					Proteins: 35.0,
					Fats:     25.0,
					Carbs:    70.0,
					Servings: 1,
				},
			},
		},
	}
	
	adults := 2
	children := 1
	
	totalP, totalF, totalC := optimizer.calculateWeeklyMacros(weeklyMenu, adults, children)
	
	// Проверяем, что макронутриенты пересчитаны с учетом порций
	expectedServings := float64(adults) + float64(children)*0.7 // 2.7
	
	expectedP := (20.0 + 40.0 + 35.0) * expectedServings
	expectedF := (15.0 + 30.0 + 25.0) * expectedServings
	expectedC := (60.0 + 80.0 + 70.0) * expectedServings
	
	if totalP != expectedP {
		t.Errorf("Ожидалось totalP = %f, получено %f", expectedP, totalP)
	}
	if totalF != expectedF {
		t.Errorf("Ожидалось totalF = %f, получено %f", expectedF, totalF)
	}
	if totalC != expectedC {
		t.Errorf("Ожидалось totalC = %f, получено %f", expectedC, totalC)
	}
}

func TestMenuOptimizer_CalculateReplacementScore(t *testing.T) {
	optimizer := NewMenuOptimizer()
	
	recipe := &models.Recipe{
		Proteins: 30.0,
		Fats:     20.0,
		Carbs:    50.0,
	}
	
	// Тест 1: Нужно больше белка
	deviationP := -0.10 // Недостаток белка
	deviationF := 0.05
	deviationC := 0.05
	
	score := optimizer.calculateReplacementScore(recipe, deviationP, deviationF, deviationC)
	
	if score <= 0 {
		t.Errorf("Ожидался положительный score для рецепта с высоким содержанием белка при недостатке белка")
	}
	
	// Тест 2: Нужно меньше белка
	deviationP = 0.10 // Избыток белка
	score = optimizer.calculateReplacementScore(recipe, deviationP, deviationF, deviationC)
	
	if score >= 0 {
		t.Errorf("Ожидался отрицательный score для рецепта с высоким содержанием белка при избытке белка")
	}
}

