package domain

import "time"

type Task struct {
    ID         int       `db:"id"`
    LevelID    int       `db:"level_id"`
    Type       string    `db:"type"` // input, choice_one, etc.
    Question   string    `db:"question"`
    Answer     *string   `db:"correct_answer"`
    XPReward   int       `db:"xp_reward"`
    OrderIndex int       `db:"order_index"`
    CreatedAt  time.Time `db:"created_at"`
    UpdatedAt  time.Time `db:"updated_at"`
}