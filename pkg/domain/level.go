package domain

import "time"

type Level struct {
	ID          int       `db:"id"`
	CourseID    int       `db:"course_id"`
	Title       string    `db:"title"`
	Description *string   `db:"description"`
	OrderIndex  int       `db:"order_index"`
	XPReward    int       `db:"xp_reward"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
