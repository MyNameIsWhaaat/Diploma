package service

import (
	"fmt"

	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
)

func (s *Service) GetByCourse(courseID int) ([]domain.LevelFromCourse, error) {

	if courseID <= 0 {
		return nil, fmt.Errorf("invalid course ID: %d", courseID)
	}

	return s.repo.GetByCourse(courseID)
}

func (s *Service) GetCourseLevelsForUser(userID, courseID int) ([]domain.LevelWithUserProgress, error) {
	if courseID <= 0 {
		return nil, fmt.Errorf("invalid course ID: %d", courseID)
	}
	if userID <= 0 {
		return nil, fmt.Errorf("invalid course ID: %d", userID)
	}
	return s.repo.GetCourseLevelsForUser(userID, courseID)
}

func (s *Service) GetLevelsWithAccess(userID, courseID int) ([]domain.LevelWithAccess, error) {
	// 1. Получаем уровни курса (упорядоченные по номеру)
	levels, err := s.repo.GetByCourse(courseID)
	if err != nil {
		return nil, fmt.Errorf("get levels: %w", err)
	}

	// 2. Получаем ID завершённых уровней пользователя
	completedIDs, err := s.repo.GetCompletedLevelIDs(userID)
	if err != nil {
		return nil, fmt.Errorf("get completed levels: %w", err)
	}
	completedMap := make(map[int]bool)
	for _, id := range completedIDs {
		completedMap[id] = true
	}

	// 3. Формируем финальный список с флагами
	var result []domain.LevelWithAccess
	for i, level := range levels {
		completed := completedMap[level.ID]
		unlocked := false

		if i == 0 {
			unlocked = true // первый уровень всегда открыт
		} else {
			prevLevel := levels[i-1]
			if completedMap[prevLevel.ID] {
				unlocked = true // предыдущий завершён — текущий открыт
			}
		}

		result = append(result, domain.LevelWithAccess{
			ID:          level.ID,
			Title:       level.Title,
			Number:      level.OrderIndex,
			IsUnlocked:  unlocked,
			IsCompleted: completed,
		})
	}

	return result, nil
}

func (s *Service) GetTheoryByLevel(levelID int) ([]domain.TheoryBlock, error) {
	if levelID <= 0 {
		return nil, fmt.Errorf("invalid level ID: %d", levelID)
	}
	return s.repo.GetTheoryBlocksByLevel(levelID)
}
