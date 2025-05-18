package domain

import "time"

type UserLevel struct {
	ID          int        `db:"id"`
	UserID      int        `db:"user_id"`
	LevelID     int        `db:"level_id"`
	IsCompleted bool       `db:"completed"`
	IsCurrent   bool       `db:"is_current"`
	XPEarned    int        `db:"xp_earned"`
	StartedAt   time.Time  `db:"started_at"`
	CompletedAt *time.Time `db:"completed_at"` // если завершён
}
