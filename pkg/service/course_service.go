package service

import (
	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
	"github.com/MyNameIsWhaaat/algo-learning/pkg/repository"
)

type CourseService struct {
	repo repository.Course
}

func NewCourseService(repo repository.Course) *CourseService {
	return &CourseService{repo: repo}
}

func (s *CourseService) GetCourseWithProgress(userId int) ([]domain.CourseWithProgress, error) {
	return s.repo.GetCourseWithProgress(userId)
}

func (s *CourseService) GetCourseWithoutProgress(userId int) ([]domain.CourseWithProgress, error) {
	return s.repo.GetCourseWithoutProgress(userId)
}

func (s *CourseService) StartCourse(userID, courseID int) error {
	return s.repo.StartCourse(userID, courseID)
}