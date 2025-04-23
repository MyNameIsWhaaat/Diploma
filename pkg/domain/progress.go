package domain

import "time"

type Progress struct {
	ID          int        `db:"id"`
	UserID      int        `db:"user_id"`
	TaskID      int        `db:"task_id"`
	IsCompleted bool       `db:"is_completed"`
	LastAnswer  *string    `db:"last_answer"`
	Attempts    int        `db:"attempts"`
	XPEarned    int        `db:"xp_earned"`
	TimeSpent   *int       `db:"time_spent"` // в секундах
	CompletedAt *time.Time `db:"completed_at"`
	NeedsReview bool       `db:"needs_review"`
}