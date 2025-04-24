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

func (r *CoursePostgres) GetAll() ([]models.CourseWithProgress, error) {
	query := `
		SELECT
			id,
			title,
			short_description,
			image_url,
			0 AS xp_earned,
			xp_reward AS xp_total,
			false AS completed
		FROM courses
	`

	var courses []models.CourseWithProgress
	err := r.db.Select(&courses, query)
	if err != nil {
		return nil, err
	}

	return courses, nil
}
