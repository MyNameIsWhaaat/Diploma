package repository

import (
	"fmt"

	"github.com/MyNameIsWhaaat/algo-learning/pkg/models"
	"github.com/jmoiron/sqlx"
)

type LevelPostgres struct {
	db *sqlx.DB
}

func NewLevelPostgres(db *sqlx.DB) *LevelPostgres {
	return &LevelPostgres{db: db}
}

func (r *LevelPostgres) GetByCourse(courseID int) ([]models.Level, error) {
	query := `
		SELECT id, course_id, title, description, order_index, xp_reward
		FROM levels
		WHERE course_id = $1
		ORDER BY order_index
	`
	var levels []models.Level
	err := r.db.Select(&levels, query, courseID)
	return levels, err
}

func (r *LevelPostgres) CompleteLevel(userID, levelID int) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Проверка: существует ли уровень
	var courseID int
	var xpReward int
	err = tx.QueryRow(`
		SELECT course_id, xp_reward FROM levels WHERE id = $1
	`, levelID).Scan(&courseID, &xpReward)
	if err != nil {
		return fmt.Errorf("уровень не найден")
	}

	// Помечаем уровень как пройденный
	_, err = tx.Exec(`
		INSERT INTO user_levels (user_id, level_id, completed)
		VALUES ($1, $2, true)
		ON CONFLICT (user_id, level_id) DO UPDATE
		SET completed = true
	`, userID, levelID)
	if err != nil {
		return err
	}

	// Обновляем XP у пользователя в курсе
	_, err = tx.Exec(`
		UPDATE user_courses
		SET xp_earned = xp_earned + $1
		WHERE user_id = $2 AND course_id = $3
	`, xpReward, userID, courseID)
	if err != nil {
		return err
	}

	// Обновляем общий XP пользователя в profile_levels
	var totalXP, xpToNextLevel int
	err = tx.QueryRow(`
		SELECT total_xp, xp_to_next_level FROM user_profile_levels WHERE user_id = $1
	`, userID).Scan(&totalXP, &xpToNextLevel)
	if err != nil {
		return err
	}

	// Прибавляем полученный опыт
	newTotalXP := totalXP + xpReward

	// Проверяем, апнулся ли уровень
	if newTotalXP >= xpToNextLevel {
		// Повышаем уровень
		_, err = tx.Exec(`
			UPDATE user_profile_levels
			SET 
				profile_level_id = profile_level_id + 1,  -- Увеличиваем ID уровня
				total_xp = $1,
				xp_to_next_level = xp_to_next_level * 2, -- Увеличиваем XP для следующего уровня
				last_level_up = NOW(),
				updated_at = NOW()
			WHERE user_id = $2
		`, newTotalXP, userID)
	} else {
		// Просто обновляем общий опыт без повышения уровня
		_, err = tx.Exec(`
			UPDATE user_profile_levels
			SET 
				total_xp = $1,
				updated_at = NOW()
			WHERE user_id = $2
		`, newTotalXP, userID)
	}

	return err
}