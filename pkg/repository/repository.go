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

type Repository struct{
	Authorization
	Course 
}

func NewRepository(db *sqlx.DB)  *Repository{
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Course: NewCoursePostgres(db),
	}
}