package services

import (
	"github.com/myplate/backend/internal/models"
	"github.com/myplate/backend/internal/repositories"
)

type ShoppingListService struct {
	shoppingRepo *repositories.ShoppingListRepository
	menuRepo     *repositories.MenuRepository
	recipeRepo   *repositories.RecipeRepository
	pantryRepo   *repositories.PantryRepository
}

func NewShoppingListService(
	shoppingRepo *repositories.ShoppingListRepository,
	menuRepo *repositories.MenuRepository,
	recipeRepo *repositories.RecipeRepository,
	pantryRepo *repositories.PantryRepository,
) *ShoppingListService {
	return &ShoppingListService{
		shoppingRepo: shoppingRepo,
		menuRepo:     menuRepo,
		recipeRepo:   recipeRepo,
		pantryRepo:   pantryRepo,
	}
}

func (s *ShoppingListService) GetByMenuID(menuID int) (*models.ShoppingList, error) {
	return s.shoppingRepo.GetByMenuID(menuID)
}


