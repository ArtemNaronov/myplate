package services

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/myplate/backend/internal/models"
	"github.com/myplate/backend/internal/repositories"
	"github.com/myplate/backend/pkg/database"
)

type MenuService struct {
	recipeRepo  *repositories.RecipeRepository
	menuRepo    *repositories.MenuRepository
	pantryRepo  *repositories.PantryRepository
	shoppingRepo *repositories.ShoppingListRepository
	goalsRepo   *repositories.GoalsRepository
}

func NewMenuService(
	recipeRepo *repositories.RecipeRepository,
	menuRepo *repositories.MenuRepository,
	pantryRepo *repositories.PantryRepository,
	shoppingRepo *repositories.ShoppingListRepository,
	goalsRepo *repositories.GoalsRepository,
) *MenuService {
	return &MenuService{
		recipeRepo:  recipeRepo,
		menuRepo:    menuRepo,
		pantryRepo:  pantryRepo,
		shoppingRepo: shoppingRepo,
		goalsRepo:   goalsRepo,
	}
}

func (s *MenuService) GenerateMenu(req *models.MenuGenerateRequest) (*models.Menu, error) {
	// Если целевые калории не указаны, берем из целей пользователя
	if req.TargetCalories == 0 {
		goals, err := s.goalsRepo.GetByUserID(req.UserID)
		if err != nil {
			return nil, fmt.Errorf("ошибка при получении целей пользователя: %w", err)
		}
		if goals != nil && goals.DailyCalories > 0 {
			req.TargetCalories = goals.DailyCalories
		} else {
			// Значение по умолчанию
			req.TargetCalories = 2000
		}
	}
	
	// Get filtered recipes
	mealTypes := []string{"breakfast", "lunch", "dinner"}
	var maxCalories, maxTime *int
	
	if req.MaxTimePerMeal > 0 {
		maxTime = &req.MaxTimePerMeal
	}
	
	recipes, err := s.recipeRepo.GetFiltered(req.DietType, req.Allergies, mealTypes, maxCalories, nil, maxTime)
	if err != nil {
		return nil, err
	}
	
	// Get pantry items if needed
	var pantryItems []models.PantryItem
	if req.ConsiderPantry {
		pantryItems, err = s.pantryRepo.GetByUserID(req.UserID)
		if err != nil {
			return nil, err
		}
	}
	
	// Score and filter recipes based on pantry
	scoredRecipes := s.scoreRecipesByPantry(recipes, pantryItems, req.PantryImportance)
	
	// Generate menu combinations
	bestMenu := s.findBestMenuCombination(scoredRecipes, req, pantryItems)
	
	if bestMenu == nil {
		return nil, fmt.Errorf("не найдена подходящая комбинация меню")
	}
	
	// Calculate totals
	totalCalories := 0
	totalTime := 0
	for _, meal := range bestMenu.Meals {
		totalCalories += meal.Calories
		totalTime += meal.Time
	}
	
	bestMenu.TotalCalories = totalCalories
	bestMenu.TotalTime = totalTime
	bestMenu.UserID = req.UserID
	bestMenu.Date = time.Now()
	
	// Calculate ingredients used and missing
	if req.ConsiderPantry {
		adults := req.Adults
		if adults == 0 {
			adults = 1 // По умолчанию 1 взрослый
		}
		children := req.Children
		bestMenu.IngredientsUsed, bestMenu.MissingIngredients = s.calculateIngredientUsage(bestMenu.Meals, recipes, pantryItems, adults, children)
	}
	
	// Save menu
	err = s.menuRepo.Create(bestMenu)
	if err != nil {
		return nil, fmt.Errorf("ошибка при сохранении меню: %w", err)
	}
	
	// Generate shopping list
	adults := req.Adults
	if adults == 0 {
		adults = 1
	}
	children := req.Children
	shoppingList := s.generateShoppingList(bestMenu, recipes, pantryItems, adults, children)
	shoppingList.UserID = req.UserID
	shoppingList.MenuID = bestMenu.ID
	err = s.shoppingRepo.CreateOrUpdate(shoppingList)
	if err != nil {
		// Логируем ошибку, но не прерываем выполнение - меню уже сохранено
		fmt.Printf("Предупреждение: не удалось сохранить список покупок для меню %d: %v\n", bestMenu.ID, err)
	}
	
	return bestMenu, nil
}

// GenerateWeeklyMenu генерирует меню на неделю (7 дней) с учетом анти-повторов и баланса БЖУ
func (s *MenuService) GenerateWeeklyMenu(req *models.WeeklyMenuRequest) (*models.WeeklyMenu, error) {
	// Рассчитываем целевые калории на день с учетом количества людей
	adults := req.Adults
	if adults == 0 {
		adults = 1
	}
	children := req.Children
	
	// Калорийная цель дня: adults*2000 + children*1400
	targetDayCalories := adults*2000 + children*1400
	
	// Распределение калорий: завтрак 25%, обед 40%, ужин 35%
	targetBreakfastCalories := int(float64(targetDayCalories) * 0.25)
	targetLunchCalories := int(float64(targetDayCalories) * 0.40)
	targetDinnerCalories := int(float64(targetDayCalories) * 0.35)
	
	// Получаем рецепты по категориям
	var maxTime *int
	if req.MaxTimePerMeal > 0 {
		maxTime = &req.MaxTimePerMeal
	}
	
	breakfastRecipes, err := s.recipeRepo.GetFiltered(req.DietType, req.Allergies, []string{"breakfast"}, nil, nil, maxTime)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении рецептов для завтрака: %w", err)
	}
	
	lunchRecipes, err := s.recipeRepo.GetFiltered(req.DietType, req.Allergies, []string{"lunch"}, nil, nil, maxTime)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении рецептов для обеда: %w", err)
	}
	
	dinnerRecipes, err := s.recipeRepo.GetFiltered(req.DietType, req.Allergies, []string{"dinner"}, nil, nil, maxTime)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении рецептов для ужина: %w", err)
	}
	
	// Получаем ингредиенты из кладовой
	var pantryItems []models.PantryItem
	if req.ConsiderPantry {
		pantryItems, err = s.pantryRepo.GetByUserID(req.UserID)
		if err != nil {
			return nil, err
		}
	}
	
	// Оцениваем рецепты
	scoredBreakfast := s.scoreRecipesByPantry(breakfastRecipes, pantryItems, req.PantryImportance)
	scoredLunch := s.scoreRecipesByPantry(lunchRecipes, pantryItems, req.PantryImportance)
	scoredDinner := s.scoreRecipesByPantry(dinnerRecipes, pantryItems, req.PantryImportance)
	
	// Проверяем, что есть хотя бы один рецепт в каждой категории
	if len(breakfastRecipes) == 0 {
		return nil, fmt.Errorf("не найдено рецептов для завтрака")
	}
	if len(lunchRecipes) == 0 {
		return nil, fmt.Errorf("не найдено рецептов для обеда")
	}
	if len(dinnerRecipes) == 0 {
		return nil, fmt.Errorf("не найдено рецептов для ужина")
	}
	
	// Генерируем меню на каждый день недели с анти-повторами
	weeklyMenu := &models.WeeklyMenu{
		Week: make([]models.WeeklyDayMenu, 7),
	}
	
	// Карта использованных рецептов для анти-повторов (не использовать 3 дня подряд)
	usedRecipeIDs := make(map[int][]int) // day -> []recipeIDs
	
	// Сохраняем ссылки на исходные списки для гарантированного fallback
	// Это важно, так как мы будем использовать их в финальном fallback
	_ = breakfastRecipes // Используем для fallback
	_ = lunchRecipes     // Используем для fallback
	_ = dinnerRecipes    // Используем для fallback
	
	for day := 0; day < 7; day++ {
		// Получаем список рецептов, которые нельзя использовать (последние 3 дня)
		excludedIDs := make(map[int]bool)
		for d := day - 3; d < day; d++ {
			if d >= 0 && usedRecipeIDs[d] != nil {
				for _, id := range usedRecipeIDs[d] {
					excludedIDs[id] = true
				}
			}
		}
		
		// Выбираем рецепты для каждого приема пищи
		var breakfastRecipe, lunchRecipe, dinnerRecipe *models.Recipe
		
		// Попытка 1: С учетом anti-repeat и калорий (если есть scored рецепты)
		if len(scoredBreakfast) > 0 {
			breakfastRecipe = s.selectRecipeForMeal(scoredBreakfast, targetBreakfastCalories, excludedIDs)
		}
		if len(scoredLunch) > 0 {
			lunchRecipe = s.selectRecipeForMeal(scoredLunch, targetLunchCalories, excludedIDs)
		}
		if len(scoredDinner) > 0 {
			dinnerRecipe = s.selectRecipeForMeal(scoredDinner, targetDinnerCalories, excludedIDs)
		}
		
		// Попытка 2: Без anti-repeat (разрешаем повторы)
		emptyExcludedIDs := make(map[int]bool)
		if breakfastRecipe == nil && len(scoredBreakfast) > 0 {
			breakfastRecipe = s.selectRecipeForMeal(scoredBreakfast, targetBreakfastCalories, emptyExcludedIDs)
		}
		if lunchRecipe == nil && len(scoredLunch) > 0 {
			lunchRecipe = s.selectRecipeForMeal(scoredLunch, targetLunchCalories, emptyExcludedIDs)
		}
		if dinnerRecipe == nil && len(scoredDinner) > 0 {
			dinnerRecipe = s.selectRecipeForMeal(scoredDinner, targetDinnerCalories, emptyExcludedIDs)
		}
		
		// Попытка 3: Игнорируем калории
		if breakfastRecipe == nil && len(scoredBreakfast) > 0 {
			breakfastRecipe = s.selectRecipeForMealIgnoreCalories(scoredBreakfast, emptyExcludedIDs)
		}
		if lunchRecipe == nil && len(scoredLunch) > 0 {
			lunchRecipe = s.selectRecipeForMealIgnoreCalories(scoredLunch, emptyExcludedIDs)
		}
		if dinnerRecipe == nil && len(scoredDinner) > 0 {
			dinnerRecipe = s.selectRecipeForMealIgnoreCalories(scoredDinner, emptyExcludedIDs)
		}
		
		// ФИНАЛЬНЫЙ FALLBACK: ВСЕГДА используем исходные списки, если рецепт не найден
		// Это гарантирует, что рецепт будет найден, так как мы проверили наличие рецептов выше
		if breakfastRecipe == nil {
			idx := day % len(breakfastRecipes)
			breakfastRecipe = &breakfastRecipes[idx]
		}
		if lunchRecipe == nil {
			idx := day % len(lunchRecipes)
			lunchRecipe = &lunchRecipes[idx]
		}
		if dinnerRecipe == nil {
			idx := day % len(dinnerRecipes)
			dinnerRecipe = &dinnerRecipes[idx]
		}
		
		// После финального fallback рецепты ОБЯЗАТЕЛЬНО должны быть найдены
		// Если это не так, значит списки пустые (но мы проверили это выше)
		if breakfastRecipe == nil || lunchRecipe == nil || dinnerRecipe == nil {
			return nil, fmt.Errorf("не найдены подходящие рецепты для дня %d (завтрак: %d рецептов, обед: %d, ужин: %d)", 
				day+1, len(breakfastRecipes), len(lunchRecipes), len(dinnerRecipes))
		}
		
		// Сохраняем использованные рецепты
		usedRecipeIDs[day] = []int{breakfastRecipe.ID, lunchRecipe.ID, dinnerRecipe.ID}
		
		// Рассчитываем итоги дня с учетом количества людей
		totalServings := float64(adults) + float64(children)*0.7
		if totalServings == 0 {
			totalServings = 1.0
		}
		
		// Пересчитываем калории и макронутриенты с учетом порций
		breakfastMultiplier := totalServings / float64(breakfastRecipe.Servings)
		if breakfastMultiplier == 0 {
			breakfastMultiplier = 1.0
		}
		lunchMultiplier := totalServings / float64(lunchRecipe.Servings)
		if lunchMultiplier == 0 {
			lunchMultiplier = 1.0
		}
		dinnerMultiplier := totalServings / float64(dinnerRecipe.Servings)
		if dinnerMultiplier == 0 {
			dinnerMultiplier = 1.0
		}
		
		totalCalories := int(float64(breakfastRecipe.Calories)*breakfastMultiplier) +
			int(float64(lunchRecipe.Calories)*lunchMultiplier) +
			int(float64(dinnerRecipe.Calories)*dinnerMultiplier)
		
		totalProteins := breakfastRecipe.Proteins*breakfastMultiplier +
			lunchRecipe.Proteins*lunchMultiplier +
			dinnerRecipe.Proteins*dinnerMultiplier
		
		totalFats := breakfastRecipe.Fats*breakfastMultiplier +
			lunchRecipe.Fats*lunchMultiplier +
			dinnerRecipe.Fats*dinnerMultiplier
		
		totalCarbs := breakfastRecipe.Carbs*breakfastMultiplier +
			lunchRecipe.Carbs*lunchMultiplier +
			dinnerRecipe.Carbs*dinnerMultiplier
		
		totalTime := breakfastRecipe.CookingTime + lunchRecipe.CookingTime + dinnerRecipe.CookingTime
		
		// Рассчитываем ингредиенты
		var ingredientsUsed, missingIngredients models.Ingredients
		if req.ConsiderPantry {
			meals := models.MenuMeals{
				{RecipeID: breakfastRecipe.ID, MealType: "breakfast", Calories: breakfastRecipe.Calories, Time: breakfastRecipe.CookingTime},
				{RecipeID: lunchRecipe.ID, MealType: "lunch", Calories: lunchRecipe.Calories, Time: lunchRecipe.CookingTime},
				{RecipeID: dinnerRecipe.ID, MealType: "dinner", Calories: dinnerRecipe.Calories, Time: dinnerRecipe.CookingTime},
			}
			allRecipes := append(append(breakfastRecipes, lunchRecipes...), dinnerRecipes...)
			ingredientsUsed, missingIngredients = s.calculateIngredientUsage(meals, allRecipes, pantryItems, adults, children)
		}
		
		// Преобразуем рецепты в DTO
		breakfastDTO := s.recipeToDTO(breakfastRecipe)
		lunchDTO := s.recipeToDTO(lunchRecipe)
		dinnerDTO := s.recipeToDTO(dinnerRecipe)
		
		weeklyMenu.Week[day] = models.WeeklyDayMenu{
			Day:                day + 1,
			Breakfast:          breakfastDTO,
			Lunch:              lunchDTO,
			Dinner:             dinnerDTO,
			TotalCalories:      totalCalories,
			TotalProteins:      totalProteins,
			TotalFats:          totalFats,
			TotalCarbs:         totalCarbs,
			TotalTime:          totalTime,
			IngredientsUsed:    ingredientsUsed,
			MissingIngredients: missingIngredients,
		}
	}
	
	// Применяем оптимизацию баланса БЖУ
	optimizer := NewMenuOptimizer()
	allRecipes := append(append(breakfastRecipes, lunchRecipes...), dinnerRecipes...)
	err = optimizer.OptimizeWeeklyMacros(weeklyMenu, allRecipes, adults, children)
	if err != nil {
		// Логируем ошибку, но не прерываем выполнение
		fmt.Printf("Предупреждение: ошибка при оптимизации БЖУ: %v\n", err)
	}
	
	return weeklyMenu, nil
}

// SaveWeeklyMenu сохраняет недельное меню в базу данных
func (s *MenuService) SaveWeeklyMenu(userID int, weeklyMenu *models.WeeklyMenu) (*models.Menu, error) {
	// Рассчитываем общие калории и время за неделю
	totalCalories := 0
	totalTime := 0
	for _, day := range weeklyMenu.Week {
		totalCalories += day.TotalCalories
		totalTime += day.TotalTime
	}
	
	// Сохраняем недельное меню в JSON формате в поле meals
	var weeklyMealsData []map[string]interface{}
	for _, day := range weeklyMenu.Week {
		dayData := map[string]interface{}{
			"day":                day.Day,
			"breakfast":          day.Breakfast,
			"lunch":              day.Lunch,
			"dinner":             day.Dinner,
			"totalCalories":      day.TotalCalories,
			"totalProteins":      day.TotalProteins,
			"totalFats":          day.TotalFats,
			"totalCarbs":          day.TotalCarbs,
			"totalTime":          day.TotalTime,
			"ingredients_used":   day.IngredientsUsed,
			"missing_ingredients": day.MissingIngredients,
		}
		weeklyMealsData = append(weeklyMealsData, dayData)
	}
	
	mealsJSON, err := json.Marshal(weeklyMealsData)
	if err != nil {
		return nil, fmt.Errorf("ошибка при сериализации данных недели: %w", err)
	}
	
	// Создаем Menu объект с недельным меню в формате JSON
	menu := &models.Menu{
		UserID:        userID,
		Date:          time.Now(), // Дата создания недельного меню
		TotalCalories: totalCalories,
		TotalTime:     totalTime,
		MenuType:      "weekly",
		Meals:         models.MenuMeals{}, // Будет заполнено через прямой SQL
	}
	
	// Сохраняем через прямой SQL запрос, так как нужно сохранить JSON напрямую
	query := `
		INSERT INTO menus (user_id, date, total_calories, total_price, total_time, menu_type, meals, ingredients_used, missing_ingredients)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`
	
	var totalPrice float64 = 0
	var ingredientsUsedJSON, missingIngredientsJSON []byte
	// Объединяем все ингредиенты из всех дней
	var allIngredientsUsed, allMissingIngredients models.Ingredients
	for _, day := range weeklyMenu.Week {
		allIngredientsUsed = append(allIngredientsUsed, day.IngredientsUsed...)
		allMissingIngredients = append(allMissingIngredients, day.MissingIngredients...)
	}
	ingredientsUsedJSON, _ = json.Marshal(allIngredientsUsed)
	missingIngredientsJSON, _ = json.Marshal(allMissingIngredients)
	
	err = database.DB.QueryRow(query,
		menu.UserID, menu.Date, menu.TotalCalories, totalPrice, menu.TotalTime, menu.MenuType,
		mealsJSON, ingredientsUsedJSON, missingIngredientsJSON,
	).Scan(&menu.ID, &menu.CreatedAt, &menu.UpdatedAt)
	
	if err != nil {
		return nil, fmt.Errorf("ошибка при сохранении недельного меню: %w", err)
	}
	
	return menu, nil
}

// GetWeeklyMenus получает все сохраненные недельные меню пользователя
func (s *MenuService) GetWeeklyMenus(userID int) ([]models.Menu, error) {
	return s.menuRepo.GetWeeklyMenusByUserID(userID)
}

// DeleteMenu удаляет меню по ID, проверяя принадлежность пользователю
func (s *MenuService) DeleteMenu(menuID, userID int) error {
	return s.menuRepo.Delete(menuID, userID)
}

// selectRecipeForMeal выбирает рецепт для приема пищи с учетом калорий и анти-повторов
func (s *MenuService) selectRecipeForMeal(scoredRecipes []ScoredRecipe, targetCalories int, excludedIDs map[int]bool) *models.Recipe {
	// Сортируем рецепты по близости к целевым калориям
	bestRecipe := (*models.Recipe)(nil)
	bestScore := math.Inf(1)
	
	for _, sr := range scoredRecipes {
		// Пропускаем исключенные рецепты
		if excludedIDs[sr.Recipe.ID] {
			continue
		}
		
		// Вычисляем отклонение от целевых калорий
		if targetCalories > 0 {
			calDiff := math.Abs(float64(sr.Recipe.Calories - targetCalories)) / float64(targetCalories)
			
			// Комбинированный score: близость к калориям + pantry score
			score := calDiff - sr.Score*0.3 // Чем выше pantry score, тем лучше
			
			if score < bestScore {
				bestScore = score
				bestRecipe = &sr.Recipe
			}
		} else {
			// Если целевые калории не указаны, используем только pantry score
			score := -sr.Score
			if score < bestScore {
				bestScore = score
				bestRecipe = &sr.Recipe
			}
		}
	}
	
	return bestRecipe
}

// selectRecipeForMealIgnoreCalories выбирает любой доступный рецепт, игнорируя калории
func (s *MenuService) selectRecipeForMealIgnoreCalories(scoredRecipes []ScoredRecipe, excludedIDs map[int]bool) *models.Recipe {
	bestRecipe := (*models.Recipe)(nil)
	bestScore := math.Inf(-1) // Ищем максимальный pantry score
	
	for _, sr := range scoredRecipes {
		// Пропускаем исключенные рецепты
		if excludedIDs[sr.Recipe.ID] {
			continue
		}
		
		// Используем только pantry score
		score := sr.Score
		
		if score > bestScore {
			bestScore = score
			bestRecipe = &sr.Recipe
		}
	}
	
	// Если не нашли по score, берем первый доступный (даже если он в excludedIDs)
	if bestRecipe == nil {
		for _, sr := range scoredRecipes {
			return &sr.Recipe
		}
	}
	
	return bestRecipe
}

// recipeToDTO преобразует Recipe в RecipeDTO
func (s *MenuService) recipeToDTO(recipe *models.Recipe) *models.RecipeDTO {
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

type ScoredRecipe struct {
	Recipe      models.Recipe
	Score       float64
	MissingCount int
	AvailableCount int
}

func (s *MenuService) scoreRecipesByPantry(recipes []models.Recipe, pantryItems []models.PantryItem, importance string) []ScoredRecipe {
	pantryMap := make(map[string]float64)
	for _, item := range pantryItems {
		key := s.normalizeIngredientName(item.Name)
		pantryMap[key] = item.Quantity
	}
	
	var scored []ScoredRecipe
	for _, recipe := range recipes {
		available := 0
		missing := 0
		
		for _, ing := range recipe.Ingredients {
			ingName := s.normalizeIngredientName(ing.Name)
			if qty, found := pantryMap[ingName]; found && qty >= ing.Quantity {
				available++
			} else {
				missing++
			}
		}
		
		total := len(recipe.Ingredients)
		if total == 0 {
			continue
		}
		
		// Filter based on importance
		if importance == "strict" && float64(missing)/float64(total) > 0.5 {
			continue // Skip recipes missing >50% ingredients
		}
		
		score := float64(available) / float64(total)
		scored = append(scored, ScoredRecipe{
			Recipe:        recipe,
			Score:         score,
			MissingCount:  missing,
			AvailableCount: available,
		})
	}
	
	return scored
}

func (s *MenuService) normalizeIngredientName(name string) string {
	// Улучшенная нормализация: приводим к нижнему регистру и убираем лишние пробелы
	normalized := strings.ToLower(name)
	// Убираем лишние пробелы
	normalized = strings.TrimSpace(normalized)
	// Заменяем множественные пробелы на один
	for strings.Contains(normalized, "  ") {
		normalized = strings.ReplaceAll(normalized, "  ", " ")
	}
	return normalized
}

func (s *MenuService) findBestMenuCombination(scoredRecipes []ScoredRecipe, req *models.MenuGenerateRequest, pantryItems []models.PantryItem) *models.Menu {
	// Group recipes by meal type
	breakfastRecipes := []ScoredRecipe{}
	lunchRecipes := []ScoredRecipe{}
	dinnerRecipes := []ScoredRecipe{}
	
	for _, sr := range scoredRecipes {
		switch sr.Recipe.MealType {
		case "breakfast":
			breakfastRecipes = append(breakfastRecipes, sr)
		case "lunch":
			lunchRecipes = append(lunchRecipes, sr)
		case "dinner":
			dinnerRecipes = append(dinnerRecipes, sr)
		}
	}
	
	// Check if we have recipes for all meal types
	if len(breakfastRecipes) == 0 || len(lunchRecipes) == 0 || len(dinnerRecipes) == 0 {
		return nil
	}
	
	// Сортируем рецепты по пригодности (лучшие сначала)
	s.sortRecipesByFitness(breakfastRecipes, req)
	s.sortRecipesByFitness(lunchRecipes, req)
	s.sortRecipesByFitness(dinnerRecipes, req)
	
	// Параметры поиска
	bestScore := math.Inf(-1)
	var bestMenu *models.Menu
	tolerance := 0.30 // ±30% допуск по калориям
	
	// Стратегия 1: Умный поиск - пробуем лучшие комбинации
	// Берем топ-N рецептов каждого типа и пробуем все комбинации
	topN := 10
	if len(breakfastRecipes) < topN {
		topN = len(breakfastRecipes)
	}
	if len(lunchRecipes) < topN {
		topN = len(lunchRecipes)
	}
	if len(dinnerRecipes) < topN {
		topN = len(dinnerRecipes)
	}
	
	// Ограничиваем для производительности
	if topN > 15 {
		topN = 15
	}
	
	// Пробуем все комбинации топ рецептов
	// Сначала ищем меню с минимальным отклонением от целевых калорий
	bestCalDiff := math.Inf(1)
	
	for i := 0; i < topN && i < len(breakfastRecipes); i++ {
		for j := 0; j < topN && j < len(lunchRecipes); j++ {
			for k := 0; k < topN && k < len(dinnerRecipes); k++ {
				breakfast := &breakfastRecipes[i]
				lunch := &lunchRecipes[j]
				dinner := &dinnerRecipes[k]
				
				totalCal := breakfast.Recipe.Calories + lunch.Recipe.Calories + dinner.Recipe.Calories
				totalTime := breakfast.Recipe.CookingTime + lunch.Recipe.CookingTime + dinner.Recipe.CookingTime
				
				// Вычисляем отклонение от целевых калорий (приоритет #1)
				calDiff := math.Abs(float64(totalCal - req.TargetCalories)) / float64(req.TargetCalories)
				
				// Сначала ищем меню с минимальным отклонением калорий
				if calDiff < bestCalDiff {
					bestCalDiff = calDiff
				}
				
				score := s.calculateMenuScore(breakfast, lunch, dinner, req, totalCal, totalTime)
				
				// Проверяем ограничения
				timeOK := req.MaxTotalTime == 0 || totalTime <= req.MaxTotalTime
				calOK := calDiff <= tolerance
				
				// Штрафы за нарушение ограничений
				if !timeOK {
					timePenalty := float64(totalTime-req.MaxTotalTime) / float64(req.MaxTotalTime)
					if timePenalty > 0.5 {
						timePenalty = 0.5 // Максимальный штраф 50%
					}
					score *= (1.0 - timePenalty)
				}
				if !calOK {
					// Штраф за отклонение калорий - чем больше отклонение, тем больше штраф
					score *= (1.0 - calDiff*0.7) // Увеличиваем штраф до 70%
				}
				
				// Бонус за точное попадание в целевые калории
				if calOK {
					// Чем ближе к целевым калориям, тем выше бонус
					calBonus := (1.0 - calDiff) * 0.2
					score += calBonus
				}
				
				// Сохраняем лучшее меню (всегда, не только если идеальное)
				if score > bestScore {
					bestScore = score
					bestMenu = &models.Menu{
						Meals: models.MenuMeals{
							{
								RecipeID: breakfast.Recipe.ID,
								MealType: "breakfast",
								Calories: breakfast.Recipe.Calories,
								Time:     breakfast.Recipe.CookingTime,
							},
							{
								RecipeID: lunch.Recipe.ID,
								MealType: "lunch",
								Calories: lunch.Recipe.Calories,
								Time:     lunch.Recipe.CookingTime,
							},
							{
								RecipeID: dinner.Recipe.ID,
								MealType: "dinner",
								Calories: dinner.Recipe.Calories,
								Time:     dinner.Recipe.CookingTime,
							},
						},
					}
				}
			}
		}
	}
	
	// Стратегия 2: Если не нашли хорошее меню, пробуем случайные комбинации
	if bestMenu == nil || bestScore < 0.3 {
		rand.Seed(time.Now().UnixNano())
		maxAttempts := 1000
		
		for i := 0; i < maxAttempts; i++ {
			breakfast := &breakfastRecipes[rand.Intn(len(breakfastRecipes))]
			lunch := &lunchRecipes[rand.Intn(len(lunchRecipes))]
			dinner := &dinnerRecipes[rand.Intn(len(dinnerRecipes))]
			
			totalCal := breakfast.Recipe.Calories + lunch.Recipe.Calories + dinner.Recipe.Calories
			totalTime := breakfast.Recipe.CookingTime + lunch.Recipe.CookingTime + dinner.Recipe.CookingTime
			
			score := s.calculateMenuScore(breakfast, lunch, dinner, req, totalCal, totalTime)
			
			// Проверяем ограничения
			timeOK := req.MaxTotalTime == 0 || totalTime <= req.MaxTotalTime
			calDiff := math.Abs(float64(totalCal - req.TargetCalories)) / float64(req.TargetCalories)
			calOK := calDiff <= tolerance
			
			// Штрафы
			if !timeOK {
				timePenalty := float64(totalTime-req.MaxTotalTime) / float64(req.MaxTotalTime)
				if timePenalty > 0.5 {
					timePenalty = 0.5
				}
				score *= (1.0 - timePenalty)
			}
			if !calOK {
				score *= (1.0 - calDiff*0.5)
			}
			
			if score > bestScore {
				bestScore = score
				bestMenu = &models.Menu{
					Meals: models.MenuMeals{
						{
							RecipeID: breakfast.Recipe.ID,
							MealType: "breakfast",
							Calories: breakfast.Recipe.Calories,
							Time:     breakfast.Recipe.CookingTime,
						},
						{
							RecipeID: lunch.Recipe.ID,
							MealType: "lunch",
							Calories: lunch.Recipe.Calories,
							Time:     lunch.Recipe.CookingTime,
						},
						{
							RecipeID: dinner.Recipe.ID,
							MealType: "dinner",
							Calories: dinner.Recipe.Calories,
							Time:     dinner.Recipe.CookingTime,
						},
					},
				}
			}
		}
	}
	
	return bestMenu
}

// sortRecipesByFitness сортирует рецепты по пригодности для меню
func (s *MenuService) sortRecipesByFitness(recipes []ScoredRecipe, req *models.MenuGenerateRequest) {
	// Сортируем по комбинированному score: pantry score + время + калории
	for i := 0; i < len(recipes); i++ {
		for j := i + 1; j < len(recipes); j++ {
			scoreI := s.calculateRecipeFitness(&recipes[i], req)
			scoreJ := s.calculateRecipeFitness(&recipes[j], req)
			if scoreI < scoreJ {
				recipes[i], recipes[j] = recipes[j], recipes[i]
			}
		}
	}
}

// calculateRecipeFitness вычисляет пригодность отдельного рецепта
func (s *MenuService) calculateRecipeFitness(sr *ScoredRecipe, req *models.MenuGenerateRequest) float64 {
	score := sr.Score * 0.4 // Pantry score (40%)
	
	// Бонус за подходящее время (30%)
	if req.MaxTimePerMeal > 0 {
		if sr.Recipe.CookingTime <= req.MaxTimePerMeal {
			timeRatio := float64(sr.Recipe.CookingTime) / float64(req.MaxTimePerMeal)
			score += (1.0 - timeRatio) * 0.3 // Быстрее = лучше
		} else {
			score -= 0.2 // Штраф за превышение времени
		}
	}
	
	// Бонус за калории, близкие к целевым (30%)
	// Целевые калории делим на 3 приема пищи
	avgCaloriesPerMeal := float64(req.TargetCalories) / 3.0
	if avgCaloriesPerMeal > 0 {
		calDiff := math.Abs(float64(sr.Recipe.Calories) - avgCaloriesPerMeal) / avgCaloriesPerMeal
		if calDiff <= 0.5 { // В пределах 50% от целевых
			calScore := (1.0 - calDiff*2) // Нормализуем
			if calScore < 0 {
				calScore = 0
			}
			score += calScore * 0.3
		} else {
			score -= 0.1 // Штраф за слишком большое отклонение
		}
	}
	
	return score
}

func (s *MenuService) calculateMenuScore(breakfast, lunch, dinner *ScoredRecipe, req *models.MenuGenerateRequest, totalCal int, totalTime int) float64 {
	score := 0.0
	
	// 1. Calorie fit score (40%) - чем ближе к цели, тем лучше
	calDiff := math.Abs(float64(totalCal - req.TargetCalories)) / float64(req.TargetCalories)
	if calDiff > 1.0 {
		calDiff = 1.0 // Ограничиваем максимальное отклонение
	}
	calScore := (1.0 - calDiff)
	score += calScore * 0.4
	
	// 2. Time score (25%) - если указано ограничение по времени
	if req.MaxTotalTime > 0 {
		if totalTime <= req.MaxTotalTime {
			timeRatio := 1.0 - (float64(totalTime) / float64(req.MaxTotalTime))
			score += timeRatio * 0.25
		} else {
			// Штраф за превышение времени
			excessRatio := float64(totalTime-req.MaxTotalTime) / float64(req.MaxTotalTime)
			score -= math.Min(excessRatio*0.5, 0.5) // Максимальный штраф 50%
		}
	} else {
		// Бонус за быстрое приготовление (если время не критично)
		avgTime := (breakfast.Recipe.CookingTime + lunch.Recipe.CookingTime + dinner.Recipe.CookingTime) / 3.0
		if avgTime < 30 {
			score += 0.1 // Бонус за быстрое меню
		}
	}
	
	// 3. Pantry match score (20%) - использование ингредиентов из кладовой
	if req.ConsiderPantry {
		avgPantryScore := (breakfast.Score + lunch.Score + dinner.Score) / 3.0
		score += avgPantryScore * 0.2
	}
	
	// 4. Macro balance score (10%) - баланс макроэлементов
	totalProteins := breakfast.Recipe.Proteins + lunch.Recipe.Proteins + dinner.Recipe.Proteins
	totalFats := breakfast.Recipe.Fats + lunch.Recipe.Fats + dinner.Recipe.Fats
	totalCarbs := breakfast.Recipe.Carbs + lunch.Recipe.Carbs + dinner.Recipe.Carbs
	
	// Идеальное соотношение: 30% белки, 30% жиры, 40% углеводы (от калорий)
	// 1г белка = 4 ккал, 1г жира = 9 ккал, 1г углеводов = 4 ккал
	proteinCal := totalProteins * 4
	fatCal := totalFats * 9
	carbCal := totalCarbs * 4
	totalMacroCal := proteinCal + fatCal + carbCal
	
	if totalMacroCal > 0 {
		proteinRatio := proteinCal / totalMacroCal
		fatRatio := fatCal / totalMacroCal
		carbRatio := carbCal / totalMacroCal
		
		// Идеальные соотношения
		idealProtein := 0.30
		idealFat := 0.30
		idealCarb := 0.40
		
		// Вычисляем отклонение от идеала
		macroDeviation := math.Abs(proteinRatio-idealProtein) + math.Abs(fatRatio-idealFat) + math.Abs(carbRatio-idealCarb)
		macroScore := 1.0 - (macroDeviation / 1.5) // Нормализуем
		if macroScore < 0 {
			macroScore = 0
		}
		score += macroScore * 0.1
	}
	
	// 5. Variety score (5%) - разнообразие рецептов
	// Бонус за разные рецепты (не повторяющиеся)
	if breakfast.Recipe.ID != lunch.Recipe.ID && 
	   breakfast.Recipe.ID != dinner.Recipe.ID && 
	   lunch.Recipe.ID != dinner.Recipe.ID {
		score += 0.05
	}
	
	// Нормализуем score в диапазон [0, 1]
	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}
	
	return score
}

// calculateIngredientUsage рассчитывает использование ингредиентов с учетом количества людей
// adults - количество взрослых, children - количество детей (дети = 0.7 коэффициента)
func (s *MenuService) calculateIngredientUsage(meals models.MenuMeals, allRecipes []models.Recipe, pantryItems []models.PantryItem, adults int, children int) (models.Ingredients, models.Ingredients) {
	// Рассчитываем коэффициент для количества порций
	// totalServings = adults + children * 0.7
	totalServings := float64(adults) + float64(children)*0.7
	if totalServings == 0 {
		totalServings = 1.0 // По умолчанию 1 порция
	}
	
	pantryMap := make(map[string]float64)
	for _, item := range pantryItems {
		key := s.normalizeIngredientName(item.Name)
		pantryMap[key] = item.Quantity
	}
	
	recipeMap := make(map[int]models.Recipe)
	for _, recipe := range allRecipes {
		recipeMap[recipe.ID] = recipe
	}
	
	used := models.Ingredients{}
	missing := models.Ingredients{}
	
	for _, meal := range meals {
		recipe, found := recipeMap[meal.RecipeID]
		if !found {
			continue
		}
		
		// Получаем количество порций рецепта (по умолчанию 1)
		recipeServings := float64(recipe.Servings)
		if recipeServings == 0 {
			recipeServings = 1.0
		}
		
		// Коэффициент для пересчета ингредиентов
		servingMultiplier := totalServings / recipeServings
		
		for _, ing := range recipe.Ingredients {
			// Пересчитываем количество ингредиента с учетом количества людей
			adjustedQuantity := ing.Quantity * servingMultiplier
			
			ingName := s.normalizeIngredientName(ing.Name)
			if qty, found := pantryMap[ingName]; found {
				usedQty := math.Min(qty, adjustedQuantity)
				used = append(used, models.Ingredient{
					Name:     ing.Name,
					Quantity: usedQty,
					Unit:     ing.Unit,
				})
				pantryMap[ingName] -= usedQty
				
				if adjustedQuantity > usedQty {
					missing = append(missing, models.Ingredient{
						Name:     ing.Name,
						Quantity: adjustedQuantity - usedQty,
						Unit:     ing.Unit,
					})
				}
			} else {
				missing = append(missing, models.Ingredient{
					Name:     ing.Name,
					Quantity: adjustedQuantity,
					Unit:     ing.Unit,
				})
			}
		}
	}
	
	return used, missing
}

// generateShoppingList генерирует список покупок с учетом количества людей
func (s *MenuService) generateShoppingList(menu *models.Menu, allRecipes []models.Recipe, pantryItems []models.PantryItem, adults int, children int) *models.ShoppingList {
	// Рассчитываем коэффициент для количества порций
	totalServings := float64(adults) + float64(children)*0.7
	if totalServings == 0 {
		totalServings = 1.0
	}
	
	recipeMap := make(map[int]models.Recipe)
	for _, recipe := range allRecipes {
		recipeMap[recipe.ID] = recipe
	}
	
	pantryMap := make(map[string]float64)
	for _, item := range pantryItems {
		key := s.normalizeIngredientName(item.Name)
		pantryMap[key] = item.Quantity
	}
	
	shoppingMap := make(map[string]*models.ShoppingItem)
	
	for _, meal := range menu.Meals {
		recipe, found := recipeMap[meal.RecipeID]
		if !found {
			continue
		}
		
		// Получаем количество порций рецепта
		recipeServings := float64(recipe.Servings)
		if recipeServings == 0 {
			recipeServings = 1.0
		}
		
		// Коэффициент для пересчета ингредиентов
		servingMultiplier := totalServings / recipeServings
		
		for _, ing := range recipe.Ingredients {
			// Пересчитываем количество ингредиента
			adjustedQuantity := ing.Quantity * servingMultiplier
			
			ingName := s.normalizeIngredientName(ing.Name)
			available := pantryMap[ingName]
			
			if adjustedQuantity > available {
				needed := adjustedQuantity - available
				
				if item, exists := shoppingMap[ingName]; exists {
					item.Quantity += needed
					item.Reason = append(item.Reason, meal.MealType)
				} else {
					shoppingMap[ingName] = &models.ShoppingItem{
						Name:     ing.Name,
						Quantity: needed,
						Unit:     ing.Unit,
						Reason:   []string{meal.MealType},
					}
				}
			}
		}
	}
	
	items := models.ShoppingItems{}
	for _, item := range shoppingMap {
		items = append(items, *item)
	}
	
	return &models.ShoppingList{
		Items: items,
	}
}

func (s *MenuService) GetDaily(userID int, date time.Time) (*models.Menu, error) {
	return s.menuRepo.GetByUserIDAndDate(userID, date)
}

func (s *MenuService) GetAllByUserID(userID int) ([]models.Menu, error) {
	return s.menuRepo.GetAllByUserID(userID)
}

func (s *MenuService) GetByID(id int) (*models.Menu, error) {
	return s.menuRepo.GetByID(id)
}

