package service

import (
	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
)

type RepoInterface interface {
	CreateUser(user domain.User) (int, error)
	GetUser(username, password string) (domain.User, error)

	GetCourseWithProgress(userId int) ([]domain.CourseWithProgress, error)
	GetCourseWithoutProgress(userId int) ([]domain.CourseWithProgress, error)
	StartCourse(userID, courseID int) error

	GetByCourse(courseID int) ([]domain.LevelFromCourse, error)
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