package service

import (
	"errors"
	"fmt"
)

func (s *Service) StartLevel(userID, levelID int) error {
	if userID <= 0 || levelID <= 0 {
		return fmt.Errorf("invalid user or level id")
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

	// Проверяем, к какому курсу принадлежит уровень
	courseID, _, err := s.repo.GetCourseIDAndXpReward(levelID)
	if err != nil {
		return fmt.Errorf("get course for level: %w", err)
	}

	// Проверяем, записан ли пользователь на курс
	enrolled, err := s.repo.IsUserEnrolledInCourse(tx, userID, courseID)
	if err != nil {
		return fmt.Errorf("check enrollment: %w", err)
	}
	if !enrolled {
		return errors.New("user is not enrolled in this course")
	}

	// Проверка: уже стартовал?
	started, err := s.repo.IsLevelStarted(userID, levelID)
	if err != nil {
		return fmt.Errorf("check level started: %w", err)
	}
	if started {
		return errors.New("level already started")
	}

	// Создаём запись о старте уровня
	if err := s.repo.CreateUserLevel(userID, levelID); err != nil {
		return fmt.Errorf("create user_level: %w", err)
	}

	return nil
}
