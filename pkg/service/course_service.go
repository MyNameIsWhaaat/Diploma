package service

import (
	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
)

func (s *Service) GetCourseWithProgress(userId int) ([]domain.CourseWithProgress, error) {
	return s.repo.GetCourseWithProgress(userId)
}

func (s *Service) GetCourseWithoutProgress(userId int) ([]domain.CourseWithProgress, error) {
	return s.repo.GetCourseWithoutProgress(userId)
}

func (s *Service) StartCourse(userID, courseID int) error {
	return s.repo.StartCourse(userID, courseID)
}