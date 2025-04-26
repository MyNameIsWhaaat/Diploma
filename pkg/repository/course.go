package repository

import (
	"fmt"

	"github.com/MyNameIsWhaaat/algo-learning/pkg/models"
)

func (r *Repository) GetCourseWithProgress(userId int) ([]models.CourseWithProgress, error) {
	query := `
		SELECT
			c.id,
			c.title,
			c.short_description,
			c.image_url,
			uc.xp_earned,
			c.xp_reward AS xp_total,
			uc.completed
		FROM user_courses uc
		JOIN courses c ON uc.course_id = c.id
		WHERE uc.user_id = $1
		ORDER BY c.id;
	`

	var courses []models.CourseWithProgress
	err := r.db.Select(&courses, query, userId)
	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (r *Repository) GetCourseWithoutProgress(userId int) ([]models.CourseWithProgress, error) {
	query := `
		SELECT
			c.id,
			c.title,
			c.short_description,
			c.image_url,
			0 AS xp_earned,
			c.xp_reward AS xp_total,
			false AS completed
		FROM courses c
		LEFT JOIN user_courses uc
			ON uc.course_id = c.id AND uc.user_id = $1
		WHERE uc.user_id IS NULL
		ORDER BY c.id;
	`

	var courses []models.CourseWithProgress
	err := r.db.Select(&courses, query, userId)
	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (r *Repository) StartCourse(userID, courseID int) error {
	// Проверка: уже есть?
	var exists bool
	err := r.db.Get(&exists, `
		SELECT EXISTS (
			SELECT 1 FROM user_courses WHERE user_id = $1 AND course_id = $2
		)
	`, userID, courseID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("курс уже начат")
	}

	// Вставка
	_, err = r.db.Exec(`
		INSERT INTO user_courses (user_id, course_id, xp_earned, completed)
		VALUES ($1, $2, 0, false)
	`, userID, courseID)

	return err
}
