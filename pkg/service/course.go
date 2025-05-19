package service

import (
	"fmt"

	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
)

func (s *Service) GetCourseWithProgress(userId int) ([]domain.CourseWithProgress, error) {
	if userId <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", userId)
	}
	return s.repo.GetCourseWithProgress(userId)
}

func (s *Service) GetCourseWithoutProgress(userId int) ([]domain.CourseWithoutProgress, error) {
	if userId <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", userId)
	}
	return s.repo.GetCourseWithoutProgress(userId)
}

func (s *Service) StartCourse(userId, courseID int) error {
	if userId <= 0 {
		return fmt.Errorf("invalid user ID: %d", userId)
	}
	return s.repo.StartCourse(userId, courseID)
}
