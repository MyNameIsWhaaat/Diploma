package service

import (
	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
	"github.com/MyNameIsWhaaat/algo-learning/pkg/models"
	"github.com/MyNameIsWhaaat/algo-learning/pkg/repository"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Course interface {
	GetCourseWithProgress(userId int) ([]models.CourseWithProgress, error)
	GetCourseWithoutProgress(userId int) ([]models.CourseWithProgress, error)
}

type Service struct {
	Authorization
	Course
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Course:        NewCourseService(repos.Course),
	}
}