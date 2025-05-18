package domain

import (
	"database/sql"
)

type LevelWithUserProgress struct {
	ID          int           `db:"id" json:"id"`
	Title       string        `db:"title" json:"title"`
	IsCompleted sql.NullBool  `db:"is_completed"`
	IsCurrent   sql.NullBool  `db:"is_current"`
	XPEarned    sql.NullInt64 `db:"xp_earned"`
	StartedAt   sql.NullTime  `db:"started_at"`
	CompletedAt sql.NullTime  `db:"completed_at"`
	OrderIndex  int           `db:"order_index" json:"order_index"`
	IsUnlocked  bool          `json:"is_unlocked"`
}
