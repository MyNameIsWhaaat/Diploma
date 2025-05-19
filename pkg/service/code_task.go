package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
)

func (s *Service) CheckUserCode(ctx context.Context, userID, taskID int, userCode string) (*domain.CodeCheckResult, error) {
	if userID <= 0 || taskID <= 0 {
		return nil, fmt.Errorf("invalid input")
	}

	// Получаем задачу из БД
	task, err := s.repo.GetCodeTaskByTaskID(taskID)
	if err != nil {
		return nil, fmt.Errorf("get task: %w", err)
	}

	// Распаковываем тесты
	var tests []domain.TestCase
	if err := json.Unmarshal(task.Tests, &tests); err != nil {
		return nil, fmt.Errorf("invalid test format: %w", err)
	}

	// Генерация кода
	testCode := generatePythonTestCode(task.FunctionName, userCode, tests)

	// Выполнение
	output, err := executeCodeWithPiston(ctx, testCode)
	if err != nil {
		return nil, fmt.Errorf("execution error: %w", err)
	}

	var passed int
	_, err = fmt.Sscanf(output, "PASSED: %d", &passed)
	if err != nil {
		return nil, fmt.Errorf("cannot parse output: %s", output)
	}

	result := &domain.CodeCheckResult{
		IsCorrect:   passed == len(tests),
		Passed:      passed,
		Total:       len(tests),
		FailedCases: nil,
	}

	ansStr := userCode
	progress := domain.Progress{
		UserID:      userID,
		TaskID:      taskID,
		IsCompleted: true,
		IsCurrent:   result.IsCorrect,
		LastAnswer:  &ansStr,
		XPEarned:    result.Passed * 5,
		CompletedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	err = s.repo.InsertOrUpdateProgress(progress)
	if err != nil {
		return nil, fmt.Errorf("save progress: %w", err)
	}
	_ = s.repo.SaveCodeAttempt(userID, taskID, userCode, result.IsCorrect, result.Passed, result.Total)

	return result, nil
}

func generatePythonTestCode(funcName, userCode string, tests []domain.TestCase) string {
	code := userCode + "\n\npassed = 0\n"
	for _, tc := range tests {
		inJSON, _ := json.Marshal(tc.Input)
		outJSON, _ := json.Marshal(tc.Output)
		code += fmt.Sprintf("if %s(%s) == %s:\n    passed += 1\n", funcName, string(inJSON), string(outJSON))
	}
	code += "print(\"PASSED:\", passed)"
	return code
}

func executeCodeWithPiston(ctx context.Context, code string) (string, error) {
	type file struct {
		Name    string `json:"name"`
		Content string `json:"content"`
	}
	type reqBody struct {
		Language string `json:"language"`
		Version  string `json:"version"`
		Files    []file `json:"files"`
	}

	payload := reqBody{
		Language: "python3",
		Version:  "3.10.0",
		Files:    []file{{Name: "main.py", Content: code}},
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", "https://emkc.org/api/v2/piston/execute", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	run, ok := res["run"].(map[string]interface{})
	if !ok {
		return "", errors.New("no run output")
	}
	out, _ := run["output"].(string)
	return out, nil
}
