package domain

import "time"

type Course struct {
	ID               int       `db:"id"`
	Title            string    `db:"title"`
	ShortDescription string    `db:"short_description"`
	FullDescription  *string   `db:"full_description"`
	ImageURL         *string   `db:"image_url"`
	Difficulty       *string   `db:"difficulty"`
	XPReward         int       `db:"xp_reward"`
	IsActive         bool      `db:"is_active"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}