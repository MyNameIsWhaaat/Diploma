package service

import (
	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
)

func (s *Service) GetByCourse(courseID int) ([]domain.LevelFromCourse, error) {
	return s.repo.GetByCourse(courseID)
}