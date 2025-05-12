package service

import (
	"time"

	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
	"github.com/jmoiron/sqlx"
)

type RepoInterface interface {
	//Пользователи
	CreateUser(user domain.User) (int, error)
	GetUser(username, password string) (domain.User, error)
	GetUserProfileData(userID int) (domain.UserProfileData, error)

	//Курсы
	GetCourseWithProgress(userId int) ([]domain.CourseWithProgress, error)
	GetCourseWithoutProgress(userId int) ([]domain.CourseWithoutProgress, error)
	StartCourse(userID, courseID int) error
	
	//Уровни
	GetByCourse(courseID int) ([]domain.LevelFromCourse, error)
	GetCourseIDAndXpReward(levelID int) (int, int, error)
	BeginTx() (*sqlx.Tx, error)
	IsUserEnrolledInCourse(tx *sqlx.Tx, userID, courseID int) (bool, error)
	IsLevelCompletedByUser(tx *sqlx.Tx, userID, levelID int) (bool, error)
	MarkLevelAsCompleted(tx *sqlx.Tx, userID, levelID int) error
	UpdateXPInCourse(tx *sqlx.Tx, userID, courseID, xpReward int) error
	GetUserProfileLevelData(tx *sqlx.Tx, userID int) (int, int, error)
	GetProfileXPBounds(tx *sqlx.Tx, profileLevelID int) (minXP, maxXP int, err error)
	UpdateUserTotalXP(tx *sqlx.Tx, userID int, totalXP int) error
	LevelUpUser(tx *sqlx.Tx, userID int, newTotalXP int) error
	GetCourseLevelsForUser(userID, courseID int) ([]domain.LevelWithUserProgress, error)

	//Задачи
	GetTasksByLevel(levelID int) ([]domain.Task, error)
	GetTaskByID(taskID int) (domain.Task, error)
	GetTaskVariants(taskID int) ([]domain.TaskVariant, error)
	InsertOrUpdateProgress(progress domain.Progress) error
	GetReviewTasks(userID, levelID int) ([]domain.Task, error)

	//
	IsLevelStarted(userID, levelID int) (bool, error)
	CreateUserLevel(userID, levelID int) error
	UpdateUserLevelCompletion(userID, levelID, xpEarned int, completedAt time.Time) error
	GetTaskProgress(userID, taskID int) (*domain.Progress, error)

}

type Service struct {
	repo RepoInterface
}

func NewService(repo RepoInterface) *Service {
	return &Service{
		repo: repo,
	}
}
