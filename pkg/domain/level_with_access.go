package domain

type LevelWithAccess struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Number      int    `json:"number"`
	IsUnlocked  bool   `json:"is_unlocked"`
	IsCompleted bool   `json:"is_completed"`
}
