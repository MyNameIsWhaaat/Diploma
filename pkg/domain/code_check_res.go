package domain

type CodeCheckResult struct {
	IsCorrect   bool              `json:"is_correct"`
	Passed      int               `json:"passed"`
	Total       int               `json:"total"`
	FailedCases []TestCase `json:"failed_cases"`
}