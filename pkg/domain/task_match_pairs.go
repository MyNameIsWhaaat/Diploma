package domain

type TaskMatchPairs struct{
	ID int `db:"id" json:"id"`
	TaskID int `db:"task_id" json:"task_id"`
    LeftText string `db:"left_text" json:"left_text"`
    RightText string `db:"right_text" json:"right_text"`
    MatchKey string `db:"match_key" json:"match_key"`
}