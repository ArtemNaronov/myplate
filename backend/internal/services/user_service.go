package services

import (
	"github.com/myplate/backend/internal/models"
	"github.com/myplate/backend/internal/repositories"
)

type UserService struct {
	goalsRepo *repositories.GoalsRepository
}

func NewUserService(goalsRepo *repositories.GoalsRepository) *UserService {
	return &UserService{
		goalsRepo: goalsRepo,
	}
}

func (s *UserService) SetGoals(userID int, goals *models.UserGoals) error {
	goals.UserID = userID
	return s.goalsRepo.CreateOrUpdate(goals)
}

func (s *UserService) GetGoals(userID int) (*models.UserGoals, error) {
	return s.goalsRepo.GetByUserID(userID)
}


