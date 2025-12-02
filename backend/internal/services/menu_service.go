package services

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/myplate/backend/internal/models"
	"github.com/myplate/backend/internal/repositories"
)

type MenuService struct {
	recipeRepo  *repositories.RecipeRepository
	menuRepo    *repositories.MenuRepository
	pantryRepo  *repositories.PantryRepository
	shoppingRepo *repositories.ShoppingListRepository
}

func NewMenuService(
	recipeRepo *repositories.RecipeRepository,
	menuRepo *repositories.MenuRepository,
	pantryRepo *repositories.PantryRepository,
	shoppingRepo *repositories.ShoppingListRepository,
) *MenuService {
	return &MenuService{
		recipeRepo:  recipeRepo,
		menuRepo:    menuRepo,
		pantryRepo:  pantryRepo,
		shoppingRepo: shoppingRepo,
	}
}

func (s *MenuService) GenerateMenu(req *models.MenuGenerateRequest) (*models.Menu, error) {
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
		bestMenu.IngredientsUsed, bestMenu.MissingIngredients = s.calculateIngredientUsage(bestMenu.Meals, recipes, pantryItems)
	}
	
	// Save menu
	err = s.menuRepo.Create(bestMenu)
	if err != nil {
		return nil, err
	}
	
	// Generate shopping list
	shoppingList := s.generateShoppingList(bestMenu, recipes, pantryItems)
	shoppingList.UserID = req.UserID
	shoppingList.MenuID = bestMenu.ID
	s.shoppingRepo.CreateOrUpdate(shoppingList)
	
	return bestMenu, nil
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
	// Simple normalization - convert to lowercase
	// In production, you'd want more sophisticated matching
	return fmt.Sprintf("%s", name)
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
	
	// Try multiple combinations
	bestScore := math.Inf(-1)
	var bestMenu *models.Menu
	tolerance := 0.30 // ±30% - более гибкий допуск
	
	// Try random combinations - увеличиваем количество попыток
	rand.Seed(time.Now().UnixNano())
	maxAttempts := 500
	
	for i := 0; i < maxAttempts; i++ {
		var breakfast, lunch, dinner *ScoredRecipe
		
		breakfast = &breakfastRecipes[rand.Intn(len(breakfastRecipes))]
		lunch = &lunchRecipes[rand.Intn(len(lunchRecipes))]
		dinner = &dinnerRecipes[rand.Intn(len(dinnerRecipes))]
		
		totalCal := breakfast.Recipe.Calories + lunch.Recipe.Calories + dinner.Recipe.Calories
		totalTime := breakfast.Recipe.CookingTime + lunch.Recipe.CookingTime + dinner.Recipe.CookingTime
		
		// Check constraints - делаем более гибкими
		timeOK := req.MaxTotalTime == 0 || totalTime <= req.MaxTotalTime
		
		// Check calorie tolerance - более гибкий
		calDiff := math.Abs(float64(totalCal - req.TargetCalories)) / float64(req.TargetCalories)
		calOK := calDiff <= tolerance
		
		// Если все ограничения соблюдены, это идеальное меню
		if timeOK && calOK {
			score := s.calculateMenuScore(breakfast, lunch, dinner, req, totalCal, totalTime)
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
		} else if bestMenu == nil {
			// Если идеального меню нет, сохраняем лучшее из всех попыток (fallback)
			score := s.calculateMenuScore(breakfast, lunch, dinner, req, totalCal, totalTime)
			// Штрафуем за нарушение ограничений
			if !timeOK {
				score *= 0.5
			}
			if !calOK {
				score *= (1.0 - calDiff) // Чем больше отклонение, тем меньше score
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

func (s *MenuService) calculateMenuScore(breakfast, lunch, dinner *ScoredRecipe, req *models.MenuGenerateRequest, totalCal int, totalTime int) float64 {
	score := 0.0
	
	// Calorie fit score (closer to target = better)
	calDiff := math.Abs(float64(totalCal - req.TargetCalories)) / float64(req.TargetCalories)
	score += (1.0 - calDiff) * 0.6
	
	// Time score (if specified)
	if req.MaxTotalTime > 0 {
		timeRatio := 1.0 - (float64(totalTime) / float64(req.MaxTotalTime))
		if timeRatio > 0 {
			score += timeRatio * 0.3
		}
	}
	
	// Pantry match score
	if req.ConsiderPantry {
		avgPantryScore := (breakfast.Score + lunch.Score + dinner.Score) / 3.0
		score += avgPantryScore * 0.1
	}
	
	return score
}

func (s *MenuService) calculateIngredientUsage(meals models.MenuMeals, allRecipes []models.Recipe, pantryItems []models.PantryItem) (models.Ingredients, models.Ingredients) {
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
		
		for _, ing := range recipe.Ingredients {
			ingName := s.normalizeIngredientName(ing.Name)
			if qty, found := pantryMap[ingName]; found {
				usedQty := math.Min(qty, ing.Quantity)
				used = append(used, models.Ingredient{
					Name:     ing.Name,
					Quantity: usedQty,
					Unit:     ing.Unit,
				})
				pantryMap[ingName] -= usedQty
				
				if ing.Quantity > usedQty {
					missing = append(missing, models.Ingredient{
						Name:     ing.Name,
						Quantity: ing.Quantity - usedQty,
						Unit:     ing.Unit,
					})
				}
			} else {
				missing = append(missing, ing)
			}
		}
	}
	
	return used, missing
}

func (s *MenuService) generateShoppingList(menu *models.Menu, allRecipes []models.Recipe, pantryItems []models.PantryItem) *models.ShoppingList {
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
		
		for _, ing := range recipe.Ingredients {
			ingName := s.normalizeIngredientName(ing.Name)
			available := pantryMap[ingName]
			
			if ing.Quantity > available {
				needed := ing.Quantity - available
				
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

