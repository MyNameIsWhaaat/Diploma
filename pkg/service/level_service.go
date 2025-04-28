package service

import (
	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
	"github.com/MyNameIsWhaaat/algo-learning/pkg/repository"
)

type LevelService struct {
	repo       repository.Level
	courseRepo repository.Course // если хочешь валидировать принадлежность
}

func NewLevelService(repo repository.Level, courseRepo repository.Course) *LevelService {
	return &LevelService{
		repo:       repo,
		courseRepo: courseRepo,
	}
}

func (s *LevelService) GetByCourse(courseID int) ([]domain.LevelFromCourse, error) {
	// можем проверять: существует ли курс
	return s.repo.GetByCourse(courseID)
}

func (s *LevelService) CompleteLevel(userID, levelID int) error {
	return s.repo.CompleteLevel(userID, levelID)
}