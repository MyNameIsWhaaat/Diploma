package repository

import (
	"github.com/MyNameIsWhaaat/algo-learning/pkg/models"
	"github.com/jmoiron/sqlx"
)

type CoursePostgres struct {
	db *sqlx.DB
}

func NewCoursePostgres(db *sqlx.DB) *CoursePostgres {
	return &CoursePostgres{db: db}
}

func (r *CoursePostgres) GetCourseWithProgress(userId int) ([]models.CourseWithProgress, error) {
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

func (r *CoursePostgres) GetCourseWithoutProgress(userId int) ([]models.CourseWithProgress, error) {
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