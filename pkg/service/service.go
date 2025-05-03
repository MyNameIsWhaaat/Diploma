package service

import (
	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
	"github.com/jmoiron/sqlx"
)

type RepoInterface interface {
	CreateUser(user domain.User) (int, error)
	GetUser(username, password string) (domain.User, error)

	GetCourseWithProgress(userId int) ([]domain.CourseWithProgress, error)
	GetCourseWithoutProgress(userId int) ([]domain.CourseWithProgress, error)
	StartCourse(userID, courseID int) error

	GetByCourse(courseID int) ([]domain.LevelFromCourse, error)
	CompleteLevel(userID, levelID int) error

	GetCourseIDAndXpReward(levelID int) (int, int, error)
	BeginTx() (*sqlx.Tx, error)
	IsUserEnrolledInCourse(tx *sqlx.Tx, userID, courseID int) (bool, error)
	IsLevelCompletedByUser(tx *sqlx.Tx, userID, levelID int) (bool, error)
	MarkLevelAsCompleted(tx *sqlx.Tx, userID, levelID int) error
	UpdateXPInCourse(tx *sqlx.Tx, userID, courseID, xpReward int) error
	GetUserProfileLevelData(tx *sqlx.Tx, userID int) (int, int, error)
	UpdateUserProfileLevel(tx *sqlx.Tx, userID, profileLevelID, newTotalXP int) error
}

type Service struct {
	repo RepoInterface
}

func NewService(repo RepoInterface) *Service {
	return &Service{
		repo: repo,
	}
}
