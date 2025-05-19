package domain

type LevelResult struct {
	XPEarned  int `json:"xp"`
	Mistakes  int `json:"mistakes"`
	TimeSpent int `json:"time"` // по желанию, если считаешь
}
