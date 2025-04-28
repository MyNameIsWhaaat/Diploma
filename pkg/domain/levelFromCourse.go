package domain

type LevelFromCourse struct {
	ID          int     `db:"id" json:"id"`
	CourseID    int     `db:"course_id" json:"-"`
	Title       string  `db:"title" json:"title"`
	Description *string `db:"description" json:"description"`
	OrderIndex int     `db:"order_index" json:"order"`
	XPReward    int     `db:"xp_reward" json:"xp"`
}