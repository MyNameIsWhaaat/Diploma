package domain

type TaskVariant struct {
    ID         int    `db:"id"`
    TaskID     int    `db:"task_id"`
    Content    string `db:"content"`
    IsCorrect  bool   `db:"is_correct"`
    OrderIndex int    `db:"order_index"`
}