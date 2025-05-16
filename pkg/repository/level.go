package repository

import (
	"database/sql"
	"time"

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

// Получить границы текущего уровня
func (r *Repository) GetProfileXPBounds(tx *sqlx.Tx, profileLevelID int) (minXP, maxXP int, err error) {
	err = tx.QueryRow(`
		SELECT min_xp, max_xp FROM profile_levels WHERE id = $1
	`, profileLevelID).Scan(&minXP, &maxXP)
	return
}

// Обновить только XP (без апа уровня)
func (r *Repository) UpdateUserTotalXP(tx *sqlx.Tx, userID int, totalXP int) error {
	_, err := tx.Exec(`
		UPDATE user_profile_levels
		SET total_xp = $1,
			updated_at = NOW()
		WHERE user_id = $2
	`, totalXP, userID)
	return err
}

// Апнуть уровень и обновить XP
func (r *Repository) LevelUpUser(tx *sqlx.Tx, userID int, newTotalXP int) error {
	_, err := tx.Exec(`
		UPDATE user_profile_levels
		SET profile_level_id = profile_level_id + 1,
			total_xp = $1,
			last_level_up = NOW(),
			updated_at = NOW()
		WHERE user_id = $2
	`, newTotalXP, userID)
	return err
}

func (r *Repository) GetCourseLevelsForUser(userID, courseID int) ([]domain.LevelWithUserProgress, error) {
	var levels []domain.LevelWithUserProgress

	query := `
	SELECT
		l.id,
		l.title,
		ul.completed AS is_completed,
		ul.is_current,
		ul.xp_earned,
		ul.started_at,
		ul.completed_at,
		l.order_index
	FROM levels l
	LEFT JOIN user_levels ul ON ul.level_id = l.id AND ul.user_id = $1
	WHERE l.course_id = $2
	ORDER BY l.id;
	`

	err := r.db.Select(&levels, query, userID, courseID)
	return levels, err
}

func (r *Repository) IsLevelStarted(userID, levelID int) (bool, error) {
	var exists bool
	err := r.db.Get(&exists, `
		SELECT EXISTS (
			SELECT 1 FROM user_levels
			WHERE user_id = $1 AND level_id = $2
		)`, userID, levelID)
	return exists, err
}

func (r *Repository) CreateUserLevel(userID, levelID int) error {
	_, err := r.db.Exec(`
		INSERT INTO user_levels (user_id, level_id, completed, is_current, xp_earned, started_at)
		VALUES ($1, $2, false, true, 0, NOW())
	`, userID, levelID)
	return err
}

func (r *Repository) UpdateUserLevelCompletion(userID, levelID, xpEarned int, completedAt time.Time) error {
	_, err := r.db.Exec(`
		UPDATE user_levels
		SET completed = true,
		    is_current = false,
		    xp_earned = $3,
		    completed_at = $4
		WHERE user_id = $1 AND level_id = $2
	`, userID, levelID, xpEarned, completedAt)
	return err
}

func (r *Repository) GetTaskProgress(userID, taskID int) (*domain.Progress, error) {
	var progress domain.Progress
	err := r.db.Get(&progress, `
		SELECT * FROM progress
		WHERE user_id = $1 AND task_id = $2
	`, userID, taskID)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &progress, err
}

func (r *Repository) GetCompletedLevelIDs(userID int) ([]int, error) {
	query := `SELECT level_id FROM user_levels WHERE user_id = $1 AND completed_at IS NOT NULL`
	var ids []int
	err := r.db.Select(&ids, query, userID)
	return ids, err
}

func (r *Repository) GetTheoryBlocksByLevel(levelID int) ([]domain.TheoryBlock, error) {
	var blocks []domain.TheoryBlock

	query := `
	SELECT id, level_id, title, content, order_index
	FROM theory_blocks
	WHERE level_id = $1
	ORDER BY order_index;
	`

	err := r.db.Select(&blocks, query, levelID)
	return blocks, err
}