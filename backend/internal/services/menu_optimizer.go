package services

import (
	"math"

	"github.com/myplate/backend/internal/models"
)

// MenuOptimizer оптимизирует недельное меню по балансу БЖУ
type MenuOptimizer struct{}

func NewMenuOptimizer() *MenuOptimizer {
	return &MenuOptimizer{}
}

// OptimizeWeeklyMacros оптимизирует недельное меню по балансу БЖУ
func (o *MenuOptimizer) OptimizeWeeklyMacros(weeklyMenu *models.WeeklyMenu, allRecipes []models.Recipe, adults int, children int) error {
	// Рассчитываем целевые калории на неделю
	weeklyCalories := float64(adults*2000 + children*1400) * 7
	
	// Целевые макронутриенты на неделю
	// 25% белки, 30% жиры, 45% углеводы
	weeklyProteins := weeklyCalories * 0.25 / 4.0 // 1г белка = 4 ккал
	weeklyFats := weeklyCalories * 0.30 / 9.0     // 1г жира = 9 ккал
	weeklyCarbs := weeklyCalories * 0.45 / 4.0    // 1г углеводов = 4 ккал
	
	// Подсчитываем текущие макронутриенты
	totalP, totalF, totalC := o.calculateWeeklyMacros(weeklyMenu, adults, children)
	
	// Проверяем отклонения
	deviationP := (totalP / weeklyProteins) - 1.0
	deviationF := (totalF / weeklyFats) - 1.0
	deviationC := (totalC / weeklyCarbs) - 1.0
	
	// Если отклонения в пределах ±7%, оптимизация не нужна
	if math.Abs(deviationP) <= 0.07 && math.Abs(deviationF) <= 0.07 && math.Abs(deviationC) <= 0.07 {
		return nil
	}
	
	// Выполняем коррекцию (максимум 4 замены на неделю, 1 на день)
	maxReplacements := 4
	replacementsCount := 0
	
	// Создаем карту использованных рецептов для анти-повторов
	usedRecipes := make(map[int][]int) // day -> []recipeIDs
	
	for dayIdx := range weeklyMenu.Week {
		if replacementsCount >= maxReplacements {
			break
		}
		
		dayMenu := &weeklyMenu.Week[dayIdx]
		
		// Проверяем каждое блюдо дня
		meals := []struct {
			recipe *models.RecipeDTO
			mealType string
		}{
			{dayMenu.Breakfast, "breakfast"},
			{dayMenu.Lunch, "lunch"},
			{dayMenu.Dinner, "dinner"},
		}
		
		for _, meal := range meals {
			if meal.recipe == nil || replacementsCount >= maxReplacements {
				continue
			}
			
			// Находим полный рецепт
			var fullRecipe *models.Recipe
			for _, r := range allRecipes {
				if r.ID == meal.recipe.ID {
					fullRecipe = &r
					break
				}
			}
			if fullRecipe == nil {
				continue
			}
			
			// Проверяем, нужно ли заменять это блюдо
			needsReplacement := o.shouldReplaceMeal(
				fullRecipe, meal.mealType, deviationP, deviationF, deviationC,
				totalP, totalF, totalC, weeklyProteins, weeklyFats, weeklyCarbs,
			)
			
			if needsReplacement {
				// Ищем альтернативу
				alternative := o.findAlternative(
					allRecipes, meal.mealType, usedRecipes, dayIdx,
					deviationP, deviationF, deviationC,
					dayMenu.TotalCalories, fullRecipe.Calories,
				)
				
				if alternative != nil {
					// Заменяем блюдо
					o.replaceMeal(dayMenu, meal.mealType, alternative, adults, children)
					
					// Обновляем использованные рецепты
					if usedRecipes[dayIdx] == nil {
						usedRecipes[dayIdx] = []int{}
					}
					usedRecipes[dayIdx] = append(usedRecipes[dayIdx], alternative.ID)
					
					// Пересчитываем макронутриенты
					totalP, totalF, totalC = o.calculateWeeklyMacros(weeklyMenu, adults, children)
					deviationP = (totalP / weeklyProteins) - 1.0
					deviationF = (totalF / weeklyFats) - 1.0
					deviationC = (totalC / weeklyCarbs) - 1.0
					
					replacementsCount++
					
					// Если отклонения стали приемлемыми, прекращаем
					if math.Abs(deviationP) <= 0.07 && math.Abs(deviationF) <= 0.07 && math.Abs(deviationC) <= 0.07 {
						return nil
					}
				}
			}
		}
	}
	
	return nil
}

// calculateWeeklyMacros подсчитывает суммарные БЖУ за неделю
func (o *MenuOptimizer) calculateWeeklyMacros(weeklyMenu *models.WeeklyMenu, adults int, children int) (float64, float64, float64) {
	totalServings := float64(adults) + float64(children)*0.7
	if totalServings == 0 {
		totalServings = 1.0
	}
	
	var totalP, totalF, totalC float64
	
	for _, day := range weeklyMenu.Week {
		if day.Breakfast != nil {
			servingMultiplier := totalServings / float64(day.Breakfast.Servings)
			if servingMultiplier == 0 {
				servingMultiplier = 1.0
			}
			totalP += day.Breakfast.Proteins * servingMultiplier
			totalF += day.Breakfast.Fats * servingMultiplier
			totalC += day.Breakfast.Carbs * servingMultiplier
		}
		if day.Lunch != nil {
			servingMultiplier := totalServings / float64(day.Lunch.Servings)
			if servingMultiplier == 0 {
				servingMultiplier = 1.0
			}
			totalP += day.Lunch.Proteins * servingMultiplier
			totalF += day.Lunch.Fats * servingMultiplier
			totalC += day.Lunch.Carbs * servingMultiplier
		}
		if day.Dinner != nil {
			servingMultiplier := totalServings / float64(day.Dinner.Servings)
			if servingMultiplier == 0 {
				servingMultiplier = 1.0
			}
			totalP += day.Dinner.Proteins * servingMultiplier
			totalF += day.Dinner.Fats * servingMultiplier
			totalC += day.Dinner.Carbs * servingMultiplier
		}
	}
	
	return totalP, totalF, totalC
}

// shouldReplaceMeal определяет, нужно ли заменять блюдо
func (o *MenuOptimizer) shouldReplaceMeal(
	recipe *models.Recipe, mealType string,
	deviationP, deviationF, deviationC float64,
	currentP, currentF, currentC, targetP, targetF, targetC float64,
) bool {
	// Если отклонения уже в пределах нормы, заменять не нужно
	if math.Abs(deviationP) <= 0.07 && math.Abs(deviationF) <= 0.07 && math.Abs(deviationC) <= 0.07 {
		return false
	}
	
	// Проверяем вклад этого блюда в перекос
	recipeP := recipe.Proteins
	recipeF := recipe.Fats
	recipeC := recipe.Carbs
	
	// Если блюдо усиливает перекос, его нужно заменить
	if deviationP > 0.07 && recipeP > targetP/21 { // 21 = 7 дней * 3 приема пищи
		return true
	}
	if deviationP < -0.07 && recipeP < targetP/21 {
		return true
	}
	
	if deviationF > 0.07 && recipeF > targetF/21 {
		return true
	}
	if deviationF < -0.07 && recipeF < targetF/21 {
		return true
	}
	
	if deviationC > 0.07 && recipeC > targetC/21 {
		return true
	}
	if deviationC < -0.07 && recipeC < targetC/21 {
		return true
	}
	
	return false
}

// findAlternative находит альтернативный рецепт
func (o *MenuOptimizer) findAlternative(
	allRecipes []models.Recipe, mealType string,
	usedRecipes map[int][]int, currentDay int,
	deviationP, deviationF, deviationC float64,
	targetDayCalories int, currentCalories int,
) *models.Recipe {
	// Получаем список рецептов, которые нельзя использовать (последние 3 дня)
	excludedIDs := make(map[int]bool)
	for day := currentDay - 3; day < currentDay; day++ {
		if day >= 0 && usedRecipes[day] != nil {
			for _, id := range usedRecipes[day] {
				excludedIDs[id] = true
			}
		}
	}
	
	var bestRecipe *models.Recipe
	bestScore := math.Inf(-1)
	
	for _, recipe := range allRecipes {
		// Проверяем категорию
		if recipe.MealType != mealType {
			continue
		}
		
		// Проверяем анти-повторы
		if excludedIDs[recipe.ID] {
			continue
		}
		
		// Проверяем близость калорий (в пределах ±20%)
		calDiff := math.Abs(float64(recipe.Calories-currentCalories)) / float64(currentCalories)
		if calDiff > 0.2 {
			continue
		}
		
		// Вычисляем score: насколько хорошо рецепт корректирует баланс
		recipePtr := recipe // Создаем указатель на копию
		score := o.calculateReplacementScore(&recipePtr, deviationP, deviationF, deviationC)
		
		if score > bestScore {
			bestScore = score
			bestRecipe = &recipePtr
		}
	}
	
	return bestRecipe
}

// calculateReplacementScore вычисляет score для замены рецепта
func (o *MenuOptimizer) calculateReplacementScore(
	recipe *models.Recipe,
	deviationP, deviationF, deviationC float64,
) float64 {
	score := 0.0
	
	// Бонус за коррекцию отклонений
	if deviationP > 0.07 {
		// Нужно меньше белка
		score -= recipe.Proteins * 0.1
	} else if deviationP < -0.07 {
		// Нужно больше белка
		score += recipe.Proteins * 0.1
	}
	
	if deviationF > 0.07 {
		// Нужно меньше жира
		score -= recipe.Fats * 0.1
	} else if deviationF < -0.07 {
		// Нужно больше жира
		score += recipe.Fats * 0.1
	}
	
	if deviationC > 0.07 {
		// Нужно меньше углеводов
		score -= recipe.Carbs * 0.1
	} else if deviationC < -0.07 {
		// Нужно больше углеводов
		score += recipe.Carbs * 0.1
	}
	
	return score
}

// replaceMeal заменяет блюдо в дневном меню
func (o *MenuOptimizer) replaceMeal(
	dayMenu *models.WeeklyDayMenu, mealType string,
	recipe *models.Recipe, adults int, children int,
) {
	recipeDTO := o.recipeToDTO(recipe)
	
	switch mealType {
	case "breakfast":
		dayMenu.Breakfast = recipeDTO
	case "lunch":
		dayMenu.Lunch = recipeDTO
	case "dinner":
		dayMenu.Dinner = recipeDTO
	}
	
	// Пересчитываем итоги дня
	dayMenu.TotalCalories = 0
	dayMenu.TotalProteins = 0
	dayMenu.TotalFats = 0
	dayMenu.TotalCarbs = 0
	dayMenu.TotalTime = 0
	
	totalServings := float64(adults) + float64(children)*0.7
	if totalServings == 0 {
		totalServings = 1.0
	}
	
	if dayMenu.Breakfast != nil {
		multiplier := totalServings / float64(dayMenu.Breakfast.Servings)
		if multiplier == 0 {
			multiplier = 1.0
		}
		dayMenu.TotalCalories += int(float64(dayMenu.Breakfast.Calories) * multiplier)
		dayMenu.TotalProteins += dayMenu.Breakfast.Proteins * multiplier
		dayMenu.TotalFats += dayMenu.Breakfast.Fats * multiplier
		dayMenu.TotalCarbs += dayMenu.Breakfast.Carbs * multiplier
		dayMenu.TotalTime += dayMenu.Breakfast.CookingTime
	}
	
	if dayMenu.Lunch != nil {
		multiplier := totalServings / float64(dayMenu.Lunch.Servings)
		if multiplier == 0 {
			multiplier = 1.0
		}
		dayMenu.TotalCalories += int(float64(dayMenu.Lunch.Calories) * multiplier)
		dayMenu.TotalProteins += dayMenu.Lunch.Proteins * multiplier
		dayMenu.TotalFats += dayMenu.Lunch.Fats * multiplier
		dayMenu.TotalCarbs += dayMenu.Lunch.Carbs * multiplier
		dayMenu.TotalTime += dayMenu.Lunch.CookingTime
	}
	
	if dayMenu.Dinner != nil {
		multiplier := totalServings / float64(dayMenu.Dinner.Servings)
		if multiplier == 0 {
			multiplier = 1.0
		}
		dayMenu.TotalCalories += int(float64(dayMenu.Dinner.Calories) * multiplier)
		dayMenu.TotalProteins += dayMenu.Dinner.Proteins * multiplier
		dayMenu.TotalFats += dayMenu.Dinner.Fats * multiplier
		dayMenu.TotalCarbs += dayMenu.Dinner.Carbs * multiplier
		dayMenu.TotalTime += dayMenu.Dinner.CookingTime
	}
}

// recipeToDTO преобразует Recipe в RecipeDTO
func (o *MenuOptimizer) recipeToDTO(recipe *models.Recipe) *models.RecipeDTO {
	return &models.RecipeDTO{
		ID:           recipe.ID,
		Name:         recipe.Name,
		Description:  recipe.Description,
		Calories:     recipe.Calories,
		Proteins:     recipe.Proteins,
		Fats:         recipe.Fats,
		Carbs:        recipe.Carbs,
		CookingTime:  recipe.CookingTime,
		Servings:     recipe.Servings,
		MealType:     recipe.MealType,
		Ingredients:  recipe.Ingredients,
		Instructions: recipe.Instructions,
	}
}

