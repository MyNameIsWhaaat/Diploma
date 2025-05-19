package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
)

func (s *Service) CompleteLevel(userID, levelID int) (*domain.LevelResult, error) {

	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", userID)
	}

	tx, err := s.repo.BeginTx()
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Получаем course_id и награду за уровень
	courseID, _, err := s.repo.GetCourseIDAndXpReward(levelID)
	if err != nil {
		return nil, fmt.Errorf("get course and xp: %w", err)
	}

	// Проверяем, записан ли пользователь на курс
	enrolled, err := s.repo.IsUserEnrolledInCourse(tx, userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("check enrollment: %w", err)
	}
	if !enrolled {
		return nil, errors.New("user is not enrolled in the course")
	}

	// Проверка: уровень начат
	started, err := s.repo.IsLevelStarted(userID, levelID)
	if err != nil {
		return nil, fmt.Errorf("check level start: %w", err)
	}
	if !started {
		return nil, errors.New("level not started")
	}

	// Проверка: не завершён ли уже
	completed, err := s.repo.IsLevelCompletedByUser(tx, userID, levelID)
	if err != nil {
		return nil, fmt.Errorf("check level completion: %w", err)
	}
	if completed {
		return nil, errors.New("level already completed")
	}

	// Получаем задачи уровня
	tasks, err := s.repo.GetTasksByLevel(levelID)
	if err != nil {
		return nil, fmt.Errorf("get tasks: %w", err)
	}

	reviewTasks, err := s.repo.GetReviewTasks(userID, levelID)
	if err != nil {
		return nil, fmt.Errorf("get review tasks: %w", err)
	}
	if len(reviewTasks) > 0 {
		return nil, errors.New("review tasks are not completed")
	}

	// Проверка прогресса по задачам
	totalXP := 0
	XPforUI := 0
	for _, task := range tasks {
		progress, err := s.repo.GetTaskProgress(userID, task.ID)
		if err != nil {
			return nil, fmt.Errorf("get progress for task %d: %w", task.ID, err)
		}
		if progress == nil || !progress.IsCompleted {
			return nil, errors.New("not all tasks completed")
		}
		totalXP += progress.XPEarned
		XPforUI += progress.XPEarned
	}

	// Обновляем user_level
	err = s.repo.UpdateUserLevelCompletion(userID, levelID, totalXP, time.Now())
	if err != nil {
		return nil, fmt.Errorf("update user_level: %w", err)
	}

	// Обновляем XP в курсе
	err = s.repo.UpdateXPInCourse(tx, userID, courseID, totalXP)
	if err != nil {
		return nil, fmt.Errorf("update xp in course: %w", err)
	}

	// Получаем данные о профиле пользователя
	profileLevelID, currentTotalXP, err := s.repo.GetUserProfileLevelData(tx, userID)
	if err != nil {
		return nil, fmt.Errorf("get profile level data: %w", err)
	}

	// Считаем новое общее количество XP
	newTotalXP := currentTotalXP + totalXP

	_, maxXP, err := s.repo.GetProfileXPBounds(tx, profileLevelID)
	if err != nil {
		return nil, fmt.Errorf("get profile level bounds: %w", err)
	}

	if newTotalXP >= maxXP {
		if err := s.repo.LevelUpUser(tx, userID, newTotalXP); err != nil {
			return nil, fmt.Errorf("level up user: %w", err)
		}
	} else {
		if err := s.repo.UpdateUserTotalXP(tx, userID, newTotalXP); err != nil {
			return nil, fmt.Errorf("update user xp: %w", err)
		}
	}

	mistakes, err := s.repo.CountMistakes(userID, levelID)
	if err != nil {
		return nil, fmt.Errorf("count mistakes: %w", err)
	}

	return &domain.LevelResult{
		XPEarned: XPforUI,
		Mistakes: mistakes,
	}, nil
}
