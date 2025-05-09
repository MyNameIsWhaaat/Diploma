package repository

import "github.com/MyNameIsWhaaat/algo-learning/pkg/domain"

func (r *Repository) GetTaskByID(taskID int) (domain.Task, error) {
	var task domain.Task
	query := `
		SELECT id, level_id, type, question, correct_answer, xp_reward, order_index, created_at, updated_at
		FROM tasks WHERE id = $1
	`
	err := r.db.Get(&task, query, taskID)
	return task, err
}

func (r *Repository) GetTaskVariants(taskID int) ([]domain.TaskVariant, error) {
	var variants []domain.TaskVariant
	query := `
		SELECT id, task_id, content, is_correct, order_index
		FROM task_variants
		WHERE task_id = $1
		ORDER BY order_index ASC
	`
	err := r.db.Select(&variants, query, taskID)
	return variants, err
}

func (r *Repository) GetTasksByLevel(levelID int) ([]domain.Task, error) {
	var tasks []domain.Task
	query := `
		SELECT id, level_id, type, question, correct_answer, xp_reward, order_index, created_at, updated_at
		FROM tasks
		WHERE level_id = $1
		ORDER BY order_index ASC
	`
	err := r.db.Select(&tasks, query, levelID)
	return tasks, err
}

func (r *Repository) GetCorrectChoiceOneVariantID(taskID int) (int, error) {
	var variantID int
	query := `
		SELECT id FROM task_variants
		WHERE task_id = $1 AND is_correct = true
		LIMIT 1
	`
	err := r.db.Get(&variantID, query, taskID)
	return variantID, err
}

func (r *Repository) GetCorrectChoiceManyVariantIDs(taskID int) ([]int, error) {
	var ids []int
	query := `
		SELECT id FROM task_variants
		WHERE task_id = $1 AND is_correct = true
	`
	err := r.db.Select(&ids, query, taskID)
	return ids, err
}

func (r *Repository) InsertOrUpdateProgress(p domain.Progress) error {
	query := `
		INSERT INTO progress (user_id, task_id, is_completed, last_answer, xp_earned, completed_at, attempts)
		VALUES ($1, $2, $3, $4, $5, $6, 1)
		ON CONFLICT (user_id, task_id)
		DO UPDATE SET
			is_completed = EXCLUDED.is_completed,
			last_answer = EXCLUDED.last_answer,
			xp_earned = EXCLUDED.xp_earned,
			attempts = progress.attempts + 1,
			completed_at = CASE
				WHEN EXCLUDED.is_completed THEN EXCLUDED.completed_at
				ELSE progress.completed_at
			END
	`
	_, err := r.db.Exec(query,
		p.UserID,
		p.TaskID,
		p.IsCompleted,
		p.LastAnswer,
		p.XPEarned,
		p.CompletedAt.Time,
	)
	return err
}