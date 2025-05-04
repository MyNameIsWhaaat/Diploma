package domain

import "time"

type UserProfileLevel struct {
	ID              int        `db:"id"`        
	UserID          int        `db:"user_id"`     
	ProfileLevelID  int        `db:"profile_level_id"`
	TotalXP         int        `db:"total_xp"`    
	XPToNextLevel   int        `db:"xp_to_next_level"`
	LastLevelUp     *time.Time `db:"last_level_up"`
	CreatedAt       time.Time  `db:"created_at"` 
	UpdatedAt       time.Time  `db:"updated_at"`
}