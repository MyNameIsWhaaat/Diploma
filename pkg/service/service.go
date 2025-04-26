package service

import (
	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
	"github.com/MyNameIsWhaaat/algo-learning/pkg/models"
)

type RepoInterface interface {
	CreateUser(user domain.User) (int, error)
	GetUser(username, password string) (domain.User, error)

	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)

	GetCourseWithProgress(userId int) ([]models.CourseWithProgress, error)
	GetCourseWithoutProgress(userId int) ([]models.CourseWithProgress, error)
	StartCourse(userID, courseID int) error

	GetByCourse(courseID int) ([]models.Level, error)
	CompleteLevel(userID, levelID int) error
}

type Service struct {
	repo RepoInterface
}

func NewService(repo RepoInterface) *Service {
	return &Service{
		repo: repo,
	}
}
