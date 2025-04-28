package service

import (
	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
	"github.com/MyNameIsWhaaat/algo-learning/pkg/repository"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Course interface {
	GetCourseWithProgress(userId int) ([]domain.CourseWithProgress, error)
	GetCourseWithoutProgress(userId int) ([]domain.CourseWithProgress, error)
	StartCourse(userID, courseID int) error
}

type Level interface {
	GetByCourse(courseID int) ([]domain.LevelFromCourse, error)
	CompleteLevel(userID, levelID int) error
}

type Service struct {
	Authorization
	Course
	Level
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Course:        NewCourseService(repos.Course),
		Level:         NewLevelService(repos.Level, repos.Course),
	}
}