package domain

type CourseWithProgress struct {
	ID               int     `db:"id" json:"id"`
	Title            string  `db:"title" json:"title"`
	ShortDescription string  `db:"short_description" json:"short_description"`
	ImageURL         *string `db:"image_url" json:"image_url"`
	XPEarned         int     `db:"xp_earned" json:"xp_earned"`
	XPReward         int     `db:"xp_reward" json:"xp_reward"`
	IsCompleted      bool    `db:"completed" json:"completed"`
}