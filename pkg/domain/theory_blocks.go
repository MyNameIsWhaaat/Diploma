package domain

type TheoryBlock struct {
	ID         int    `db:"id" json:"id"`
	LevelID    int    `db:"level_id" json:"level_id"`
	Title      string `db:"title" json:"title"`
	Content    string `db:"content" json:"content"` // можно хранить HTML или Markdown
	OrderIndex int    `db:"order_index" json:"order_index"`
}
