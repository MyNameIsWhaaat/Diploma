package service

import (
	"github.com/MyNameIsWhaaat/algo-learning/pkg/models"
)

func (s *Service) GetByCourse(courseID int) ([]models.Level, error) {
	// можем проверять: существует ли курс
	return s.repo.GetByCourse(courseID)
}

func (s *Service) CompleteLevel(userID, levelID int) error {
	return s.repo.CompleteLevel(userID, levelID)
}
