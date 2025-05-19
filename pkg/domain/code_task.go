package domain

import "encoding/json"

type CodeTask struct {
	TaskID             int             `json:"task_id"`
	FunctionName       string          `json:"function_name"`
	Language           string          `json:"language"`
	Tests              json.RawMessage `json:"tests"`
	DefaultCode        string          `json:"default_code"`
	Hint               string          `json:"hint"`
	MaxExecutionTimeMS int             `json:"max_execution_time_ms"`
}
