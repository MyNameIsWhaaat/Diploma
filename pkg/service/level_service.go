package service

import (
	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
)

func (s *Service) GetByCourse(courseID int) ([]domain.LevelFromCourse, error) {
	return s.repo.GetByCourse(courseID)
}

func (s *Service) CompleteLevel(userID, levelID int) error {
	return s.repo.CompleteLevel(userID, levelID)
}