package services

import (
	"github.com/myplate/backend/internal/models"
	"github.com/myplate/backend/internal/repositories"
)

type PantryService struct {
	pantryRepo *repositories.PantryRepository
}

func NewPantryService(pantryRepo *repositories.PantryRepository) *PantryService {
	return &PantryService{
		pantryRepo: pantryRepo,
	}
}

func (s *PantryService) GetByUserID(userID int) ([]models.PantryItem, error) {
	return s.pantryRepo.GetByUserID(userID)
}

func (s *PantryService) Create(userID int, item *models.PantryItem) error {
	item.UserID = userID
	return s.pantryRepo.Create(item)
}

func (s *PantryService) Delete(userID int, id int) error {
	return s.pantryRepo.Delete(id, userID)
}


