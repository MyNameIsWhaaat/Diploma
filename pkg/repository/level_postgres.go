package repository

import (
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