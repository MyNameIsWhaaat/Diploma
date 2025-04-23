package domain

import "time"

type UserCourse struct {
	ID             int       `db:"id"`
	UserID         int       `db:"user_id"`
	CourseID       int       `db:"course_id"`
	XPEarned       int       `db:"xp_earned"`
	CurrentLevelID *int      `db:"current_level_id"` // может быть nil
	Completed      bool      `db:"completed"`
	StartedAt      time.Time `db:"started_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}