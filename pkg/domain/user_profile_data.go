package domain

type UserProfileData struct {
	ID           int    `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Username     string `json:"username" db:"username"`
	Email        string `json:"email" db:"email"`
	AvatarURL    string `json:"avatar_url" db:"avatar_url"`
	ProfileLevel string `json:"profile_level" db:"profile_level"`
	TotalXP      int    `json:"total_xp" db:"total_xp"`
	IsNewUser 	 bool 	`db:"is_new_user"`
}