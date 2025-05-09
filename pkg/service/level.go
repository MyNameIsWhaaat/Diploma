package service

import (
	"fmt"

	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
)

func (s *Service) GetByCourse(courseID int) ([]domain.LevelFromCourse, error) {

	if courseID <= 0 {
		return nil, fmt.Errorf("invalid course ID: %d", courseID)
	}

	return s.repo.GetByCourse(courseID)
}

func (s *Service) GetCourseLevelsForUser(userID, courseID int) ([]domain.LevelWithUserProgress, error){
	if courseID <= 0 {
		return nil, fmt.Errorf("invalid course ID: %d", courseID)
	}
	if userID <= 0 {
		return nil, fmt.Errorf("invalid course ID: %d", userID)
	}
	return s.repo.GetCourseLevelsForUser(userID, courseID)
}