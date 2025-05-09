package domain

type CourseWithoutProgress struct {
	ID               int     `db:"id" json:"id"`
	Title            string  `db:"title" json:"title"`
	ShortDescription string  `db:"short_description" json:"short_description"`
	FullDescription  string   `db:"full_description"`
	ImageURL         *string `db:"image_url" json:"image_url"`
	XpReward         int     `db:"xp_reward" json:"xp_reward"`
}