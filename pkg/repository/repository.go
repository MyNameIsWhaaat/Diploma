package repository

import (
	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
	"github.com/MyNameIsWhaaat/algo-learning/pkg/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface{
	CreateUser(user domain.User) (int, error)
	GetUser(username, password string) (domain.User, error)
}

type Course interface {
	GetCourseWithProgress(userId int) ([]models.CourseWithProgress, error)
	GetCourseWithoutProgress(userId int) ([]models.CourseWithProgress, error)
	StartCourse(userID, courseID int) error
}

type Level interface {
	GetByCourse(courseID int) ([]models.Level, error)
}

type Repository struct{
	Authorization
	Course
	Level
}

func NewRepository(db *sqlx.DB)  *Repository{
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Course: NewCoursePostgres(db),
		Level: NewLevelPostgres(db),
	}
}