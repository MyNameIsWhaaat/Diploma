package domain

import "time"

type UserProfileLevel struct {
	ID            int        `db:"id"`
	UserID        int        `db:"user_id"`
	CurrentLevel  int        `db:"current_level"`
	TotalXP       int        `db:"total_xp"`
	XPToNextLevel int        `db:"xp_to_next_level"`
	Title         *string    `db:"title"` // например, “Новичок”, “Гуру”
	LastLevelUp   *time.Time `db:"last_level_up"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
}