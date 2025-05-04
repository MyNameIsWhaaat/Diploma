package service

import (
	"errors"
	"fmt"
)

func (s *Service) CompleteLevel(userID, levelID int) error {

	if userID <= 0 {
		return fmt.Errorf("invalid user ID: %d", userID)
	}

	tx, err := s.repo.BeginTx()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
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
	courseID, xpReward, err := s.repo.GetCourseIDAndXpReward(levelID)
	if err != nil {
		return fmt.Errorf("get course and xp: %w", err)
	}

	// Проверяем, записан ли пользователь на курс
	enrolled, err := s.repo.IsUserEnrolledInCourse(tx, userID, courseID)
	if err != nil {
		return fmt.Errorf("check enrollment: %w", err)
	}
	if !enrolled {
		return errors.New("user is not enrolled in the course")
	}

	// Проверяем, завершил ли уже пользователь уровень
	completed, err := s.repo.IsLevelCompletedByUser(tx, userID, levelID)
	if err != nil {
		return fmt.Errorf("check level completion: %w", err)
	}
	if completed {
		return errors.New("level already completed")
	}

	// Помечаем уровень как завершённый
	err = s.repo.MarkLevelAsCompleted(tx, userID, levelID)
	if err != nil {
		return fmt.Errorf("mark level as completed: %w", err)
	}

	// Обновляем XP в курсе
	err = s.repo.UpdateXPInCourse(tx, userID, courseID, xpReward)
	if err != nil {
		return fmt.Errorf("update xp in course: %w", err)
	}

	// Получаем данные о профиле пользователя
	profileLevelID, totalXP, err := s.repo.GetUserProfileLevelData(tx, userID)
	if err != nil {
		return fmt.Errorf("get profile level data: %w", err)
	}

	// Считаем новое общее количество XP
	newTotalXP := totalXP + xpReward

	_, maxXP, err := s.repo.GetProfileXPBounds(tx, profileLevelID)
	if err != nil {
		return fmt.Errorf("get profile level bounds: %w", err)
	}

	if newTotalXP >= maxXP {
		if err := s.repo.LevelUpUser(tx, userID, newTotalXP); err != nil {
			return fmt.Errorf("level up user: %w", err)
		}
	} else {
		if err := s.repo.UpdateUserTotalXP(tx, userID, newTotalXP); err != nil {
			return fmt.Errorf("update user xp: %w", err)
		}
	}

	return nil
}