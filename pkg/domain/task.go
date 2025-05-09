package domain

import (
	"database/sql"
	"time"
)

type TaskType string

const (
	TaskTypeInput      TaskType = "input"
	TaskTypeChoiceOne  TaskType = "choice_one"
	TaskTypeChoiceMany TaskType = "choice_many"
	TaskTypeCode       TaskType = "code"
	TaskTypeMatch      TaskType = "match"
)

type Task struct {
	ID            int       `db:"id"`
	LevelID       int       `db:"level_id"`
	Type          TaskType  `db:"type"`
	Question      string    `db:"question"`
	CorrectAnswer sql.NullString    `db:"correct_answer"`
	XPReward      int       `db:"xp_reward"`
	OrderIndex    int       `db:"order_index"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}