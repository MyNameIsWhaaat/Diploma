package service

import (
	"github.com/MyNameIsWhaaat/algo-learning/pkg/models"
	"github.com/MyNameIsWhaaat/algo-learning/pkg/repository"
)

type CourseService struct {
	repo repository.Course
}

func NewCourseService(repo repository.Course) *CourseService {
	return &CourseService{repo: repo}
}

func (s *CourseService) GetAll() ([]models.CourseWithProgress, error) {
	return s.repo.GetAll()
}