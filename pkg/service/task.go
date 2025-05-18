package service

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
)

func (s *Service) SubmitAnswer(userID, taskID int, rawAnswer interface{}) (bool, int, error) {
	task, err := s.repo.GetTaskByID(taskID)
	if err != nil {
		return false, 0, fmt.Errorf("get task: %w", err)
	}

	isCorrect := false

	switch task.Type {
	case domain.TaskTypeInput:
		answerStr, ok := rawAnswer.(string)
		if !ok {
			return false, 0, errors.New("invalid input answer format")
		}
		isCorrect = normalize(answerStr) == normalize(task.CorrectAnswer.String)

	case domain.TaskTypeChoiceOne, domain.TaskTypeChoiceMany:
		variants, err := s.repo.GetTaskVariants(taskID)
		if err != nil {
			return false, 0, err
		}

		if task.Type == domain.TaskTypeChoiceOne {
			answerFloat, ok := rawAnswer.(float64)
			if !ok {
				return false, 0, errors.New("invalid choice_one answer format")
			}
			answerID := int(answerFloat)
			var correctID int
			for _, v := range variants {
				if v.IsCorrect {
					correctID = v.ID
					break
				}
			}
			isCorrect = answerID == correctID

		} else {
			answerSlice, ok := rawAnswer.([]interface{})
			if !ok {
				return false, 0, errors.New("invalid choice_many answer format")
			}
			var answers []int
			for _, a := range answerSlice {
				if f, ok := a.(float64); ok {
					answers = append(answers, int(f))
				}
			}
			var correctIDs []int
			for _, v := range variants {
				if v.IsCorrect {
					correctIDs = append(correctIDs, v.ID)
				}
			}
			isCorrect = equalIntSlices(answers, correctIDs)
		}
	case domain.TaskTypeMatch:
		// Пример ожидаемого формата: []string, где каждая строка — "ключ : значение"
		rawList, ok := rawAnswer.([]interface{})
		if !ok {
			return false, 0, errors.New("invalid match answer format")
		}

		// Получаем пары из БД
		pairs, err := s.repo.GetMatchPairsByTaskID(taskID)
		if err != nil {
			return false, 0, fmt.Errorf("get match pairs: %w", err)
		}
		for _, p := range pairs {
			fmt.Println(p.LeftText, "→", p.RightText)
		}

		// Собираем ожидаемую map[LeftText]RightText
		expected := make(map[string]string)
		for _, p := range pairs {
			expected[p.LeftText] = p.RightText // теперь сравнение по RightText, не MatchKey
		}

		// Разбираем входящие строки в map
		answerMap := make(map[string]string)
		for _, item := range rawList {
			s, ok := item.(string)
			if !ok {
				return false, 0, errors.New("match answer should be strings")
			}

			parts := strings.SplitN(s, ":", 2)
			if len(parts) != 2 {
				return false, 0, errors.New("invalid pair format (expected 'left : right')")
			}

			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])
			answerMap[left] = right
		}

		// Сравниваем
		isCorrect = true
		for key, expectedValue := range expected {
			userValue, ok := answerMap[key]
			if !ok || userValue != expectedValue {
				isCorrect = false
				break
			}
		}
	default:
		return false, 0, fmt.Errorf("unsupported task type: %s", task.Type)
	}

	// Получаем текущий прогресс, чтобы узнать попытки
	existing, err := s.repo.GetTaskProgress(userID, taskID)
	if err != nil {
		return false, 0, fmt.Errorf("get existing progress: %w", err)
	}

	xp := 0
	if isCorrect {
		if existing == nil || existing.Attempts == 0 {
			xp = task.XPReward // только за первую попытку
		}
	}

	ansStr := fmt.Sprintf("%v", rawAnswer)

	progress := domain.Progress{
		UserID:      userID,
		TaskID:      taskID,
		IsCompleted: true,
		IsCurrent:   isCorrect,
		LastAnswer:  &ansStr,
		XPEarned:    xp,
		CompletedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	err = s.repo.InsertOrUpdateProgress(progress)
	if err != nil {
		return false, 0, fmt.Errorf("save progress: %w", err)
	}

	return isCorrect, xp, nil
}

// normalize сравнивает input без учёта регистра и пробелов
func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// equalIntSlices проверяет, равны ли два списка ID (независимо от порядка)
func equalIntSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Ints(a)
	sort.Ints(b)
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func (s *Service) GetTasksByLevel(levelID int) ([]domain.Task, error) {
	if levelID <= 0 {
		return nil, fmt.Errorf("invalid level ID: %d", levelID)
	}

	return s.repo.GetTasksByLevel(levelID)
}

func (s *Service) GetTaskVariants(taskID int) ([]domain.TaskVariant, error) {
	if taskID <= 0 {
		return nil, fmt.Errorf("invalid task ID: %d", taskID)
	}

	return s.repo.GetTaskVariants(taskID)
}

func (s *Service) GetReviewTasks(userID, levelID int) ([]domain.Task, error) {
	if levelID <= 0 {
		return nil, fmt.Errorf("invalid level ID: %d", levelID)
	}
	if userID <= 0 {
		return nil, fmt.Errorf("invalid level ID: %d", userID)
	}

	return s.repo.GetReviewTasks(userID, levelID)
}

func (s *Service) CountMistakes(userID, levelID int) (int, error) {
	if levelID <= 0 {
		return 0, fmt.Errorf("invalid level ID: %d", levelID)
	}
	if userID <= 0 {
		return 0, fmt.Errorf("invalid level ID: %d", userID)
	}

	return s.repo.CountMistakes(userID, levelID)
}

func (s *Service) GetMatchPairsByTaskID(taskID int) ([]domain.TaskMatchPairs, error) {
	return s.repo.GetMatchPairsByTaskID(taskID)
}
