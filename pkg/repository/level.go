package repository

import (
	"database/sql"

	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
	"github.com/jmoiron/sqlx"
)

func (r *Repository) BeginTx() (*sqlx.Tx, error) {
	return r.db.Beginx()
}

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

func (r *Repository) GetCourseIDAndXpReward(levelID int) (int, int, error) {
	var courseID, xpReward int
	err := r.db.QueryRow(`
		SELECT course_id, xp_reward FROM levels WHERE id = $1
	`, levelID).Scan(&courseID, &xpReward)
	if err != nil {
		return 0, 0, err
	}
	return courseID, xpReward, nil
}

func (r *Repository) IsUserEnrolledInCourse(tx *sqlx.Tx, userID, courseID int) (bool, error) {
	var exists bool
	err := tx.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM user_courses
			WHERE user_id = $1 AND course_id = $2
		)
	`, userID, courseID).Scan(&exists)
	return exists, err
}

func (r *Repository) IsLevelCompletedByUser(tx *sqlx.Tx, userID, levelID int) (bool, error) {
	var completed bool
	err := tx.QueryRow(`
		SELECT completed FROM user_levels
		WHERE user_id = $1 AND level_id = $2
	`, userID, levelID).Scan(&completed)

	if err == sql.ErrNoRows {
		return false, nil
	}
	return completed, err
}

func (r *Repository) MarkLevelAsCompleted(tx *sqlx.Tx, userID, levelID int) error {
	_, err := tx.Exec(`
		INSERT INTO user_levels (user_id, level_id, completed)
		VALUES ($1, $2, true)
		ON CONFLICT (user_id, level_id) DO UPDATE
		SET completed = true
	`, userID, levelID)
	return err
}

func (r *Repository) UpdateXPInCourse(tx *sqlx.Tx, userID, courseID, xpReward int) error {
	_, err := tx.Exec(`
		INSERT INTO user_courses (user_id, course_id, xp_earned)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, course_id) DO UPDATE
		SET xp_earned = user_courses.xp_earned + EXCLUDED.xp_earned
	`, userID, courseID, xpReward)
	return err
}

func (r *Repository) GetUserProfileLevelData(tx *sqlx.Tx, userID int) (int, int, error) {
	var profileLevelID, totalXP int
	err := tx.QueryRow(`
		SELECT profile_level_id, total_xp FROM user_profile_levels WHERE user_id = $1
	`, userID).Scan(&profileLevelID, &totalXP)
	return profileLevelID, totalXP, err
}

func (r *Repository) UpdateUserProfileLevel(tx *sqlx.Tx, userID, profileLevelID, newTotalXP int) error {
	var minXP, maxXP int
	err := tx.QueryRow(`
		SELECT min_xp, max_xp FROM profile_levels WHERE id = $1
	`, profileLevelID).Scan(&minXP, &maxXP)
	if err != nil {
		return err
	}

	if newTotalXP >= maxXP {
		_, err = tx.Exec(`
			UPDATE user_profile_levels
			SET profile_level_id = profile_level_id + 1,
				total_xp = $1,
				last_level_up = NOW(),
				updated_at = NOW()
			WHERE user_id = $2
		`, newTotalXP, userID)
	} else {
		_, err = tx.Exec(`
			UPDATE user_profile_levels
			SET total_xp = $1,
				updated_at = NOW()
			WHERE user_id = $2
		`, newTotalXP, userID)
	}
	return err
}

// // Основной метод для завершения уровня
// func (r *Repository) CompleteLevel(userID, levelID int) error {
// 	tx, err := r.db.Beginx()
// 	if err != nil {
// 		return err
// 	}

// 	defer func() {
// 		if err != nil {
// 			tx.Rollback()
// 		} else {
// 			tx.Commit()
// 		}
// 	}()

// 	// Получаем данные уровня
// 	courseID, xpReward, err := r.GetCourseIDAndXpReward(levelID)
// 	if err != nil {
// 		return fmt.Errorf("уровень не найден или удалён")
// 	}

// 	//пользователь должен быть записан на курс
// 	var exists bool
// 	err = tx.QueryRow(`
// 	SELECT EXISTS(
// 		SELECT 1 FROM user_courses
// 		WHERE user_id = $1 AND course_id = $2
// 	)
// `, userID, courseID).Scan(&exists)
// 	if err != nil {
// 		return err
// 	}
// 	if !exists {
// 		return fmt.Errorf("вы не записаны на курс, содержащий этот уровень")
// 	}

// 	// Проверяем, завершён ли уже
// 	var completed bool
// 	err = tx.QueryRow(`
// 		SELECT completed FROM user_levels
// 		WHERE user_id = $1 AND level_id = $2
// 	`, userID, levelID).Scan(&completed)
// 	if err != nil && err != sql.ErrNoRows {
// 		return err
// 	}
// 	if completed {
// 		return fmt.Errorf("вы уже завершили этот уровень")
// 	}

// 	// Помечаем уровень как завершённый
// 	err = r.MarkLevelAsCompleted(tx, userID, levelID)
// 	if err != nil {
// 		return err
// 	}

// 	// Обновляем XP в курсе
// 	err = r.UpdateXPInCourse(tx, userID, courseID, xpReward)
// 	if err != nil {
// 		return err
// 	}

// 	// Получаем профиль пользователя
// 	profileLevelID, totalXP, err := r.getUserProfileLevelData(tx, userID)
// 	if err != nil {
// 		return err
// 	}

// 	// Проверяем на повышение уровня
// 	newTotalXP := totalXP + xpReward
// 	err = r.levelUp(tx, userID, profileLevelID, newTotalXP)

// 	return err
// }
