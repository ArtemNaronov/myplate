package services

import (
	"github.com/myplate/backend/internal/models"
	"github.com/myplate/backend/internal/repositories"
)

type RecipeService struct {
	recipeRepo *repositories.RecipeRepository
}

func NewRecipeService(recipeRepo *repositories.RecipeRepository) *RecipeService {
	return &RecipeService{
		recipeRepo: recipeRepo,
	}
}

func (s *RecipeService) GetAll() ([]models.Recipe, error) {
	return s.recipeRepo.GetAll()
}

func (s *RecipeService) GetByID(id int) (*models.Recipe, error) {
	return s.recipeRepo.GetByID(id)
}


