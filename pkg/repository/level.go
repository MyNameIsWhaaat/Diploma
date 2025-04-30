package repository

import (
	"database/sql"
	"fmt"

	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
	"github.com/jmoiron/sqlx"
)

func (r *Repository) GetByCourse(courseID int) ([]domain.LevelFromCourse, error) {
	query := `
		SELECT id, course_id, title, description, order_index, xp_reward
		FROM levels
		WHERE course_id = $1
		ORDER BY order_index
	`
	var levels []domain.LevelFromCourse
	err := r.db.Select(&levels, query, courseID)
	return levels, err
}

func (r *Repository) getLevelData(levelID int) (int, int, error) {
	var courseID, xpReward int
	err := r.db.QueryRow(`
		SELECT course_id, xp_reward FROM levels WHERE id = $1
	`, levelID).Scan(&courseID, &xpReward)
	if err != nil {
		return 0, 0, fmt.Errorf("уровень не найден")
	}
	return courseID, xpReward, nil
}

// Помечаем уровень как пройденный
func (r *Repository) markLevelAsCompleted(tx *sqlx.Tx, userID, levelID int) error {
	_, err := tx.Exec(`
		INSERT INTO user_levels (user_id, level_id, completed)
		VALUES ($1, $2, true)
		ON CONFLICT (user_id, level_id) DO UPDATE
		SET completed = true
	`, userID, levelID)
	return err
}

// Обновляем XP пользователя в курсе
func (r *Repository) updateXPInCourse(tx *sqlx.Tx, userID, courseID, xpReward int) error {
	// Проверяем, существует ли запись в user_courses
	var existingXP int
	err := tx.QueryRow(`
		SELECT xp_earned FROM user_courses WHERE user_id = $1 AND course_id = $2
	`, userID, courseID).Scan(&existingXP)

	// Если записи нет, то вставляем новую запись
	if err == sql.ErrNoRows {
		_, err = tx.Exec(`
			INSERT INTO user_courses (user_id, course_id, xp_earned)
			VALUES ($1, $2, $3)
		`, userID, courseID, xpReward)
		if err != nil {
			return err
		}
	} else if err != nil {
		// Если ошибка не "нет строк", то возвращаем её
		return err
	} else {
		// Если запись существует, обновляем xp_earned
		_, err = tx.Exec(`
			UPDATE user_courses
			SET xp_earned = xp_earned + $1
			WHERE user_id = $2 AND course_id = $3
		`, xpReward, userID, courseID)
		if err != nil {
			return err
		}
	}

	return nil
}

// Получаем текущий уровень пользователя и общий XP
func (r *Repository) getUserProfileLevelData(tx *sqlx.Tx, userID int) (int, int, error) {
	var profileLevelID, totalXP int
	err := tx.QueryRow(`
		SELECT profile_level_id, total_xp FROM user_profile_levels WHERE user_id = $1
	`, userID).Scan(&profileLevelID, &totalXP)
	if err != nil {
		return 0, 0, err
	}
	return profileLevelID, totalXP, nil
}

// Проверяем повышение уровня пользователя
func (r *Repository) levelUp(tx *sqlx.Tx, userID, profileLevelID, newTotalXP int) error {
	var minXP, maxXP int
	err := tx.QueryRow(`
		SELECT min_xp, max_xp FROM profile_levels WHERE id = $1
	`, profileLevelID).Scan(&minXP, &maxXP)
	if err != nil {
		return err
	}

	// Если пользователь апнулся на новый уровень
	if newTotalXP >= maxXP {
		_, err = tx.Exec(`
			UPDATE user_profile_levels
			SET 
				profile_level_id = profile_level_id + 1,  -- Увеличиваем уровень
				total_xp = $1,
				last_level_up = NOW(),
				updated_at = NOW()
			WHERE user_id = $2
		`, newTotalXP, userID)
	} else {
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

// Основной метод для завершения уровня
func (r *Repository) CompleteLevel(userID, levelID int) error {
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

	// Получаем данные уровня
	courseID, xpReward, err := r.getLevelData(levelID)
	if err != nil {
		return fmt.Errorf("уровень не найден или удалён")
	}

	//пользователь должен быть записан на курс
	var exists bool
	err = tx.QueryRow(`
	SELECT EXISTS(
		SELECT 1 FROM user_courses
		WHERE user_id = $1 AND course_id = $2
	)
`, userID, courseID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("вы не записаны на курс, содержащий этот уровень")
	}

	// Проверяем, завершён ли уже
	var completed bool
	err = tx.QueryRow(`
		SELECT completed FROM user_levels
		WHERE user_id = $1 AND level_id = $2
	`, userID, levelID).Scan(&completed)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if completed {
		return fmt.Errorf("вы уже завершили этот уровень")
	}

	// Помечаем уровень как завершённый
	err = r.markLevelAsCompleted(tx, userID, levelID)
	if err != nil {
		return err
	}

	// Обновляем XP в курсе
	err = r.updateXPInCourse(tx, userID, courseID, xpReward)
	if err != nil {
		return err
	}

	// Получаем профиль пользователя
	profileLevelID, totalXP, err := r.getUserProfileLevelData(tx, userID)
	if err != nil {
		return err
	}

	// Проверяем на повышение уровня
	newTotalXP := totalXP + xpReward
	err = r.levelUp(tx, userID, profileLevelID, newTotalXP)

	return err
}
