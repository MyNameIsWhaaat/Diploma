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
		INSERT INTO progress (user_id, task_id, is_completed, is_current, last_answer, attempts, xp_earned, time_spent, completed_at)
		VALUES (:user_id, :task_id, :is_completed, :is_current, :last_answer, 1, :xp_earned, :time_spent, :completed_at)
		ON CONFLICT (user_id, task_id) DO UPDATE
		SET 
			is_completed = EXCLUDED.is_completed,
			is_current = EXCLUDED.is_current,
			last_answer = EXCLUDED.last_answer,
			attempts = progress.attempts + 1,
			xp_earned = EXCLUDED.xp_earned,
			time_spent = EXCLUDED.time_spent,
			completed_at = EXCLUDED.completed_at
	`
	_, err := r.db.NamedExec(query, p)
	return err
}

func (r *Repository) GetReviewTasks(userID, levelID int) ([]domain.Task, error) {
	query := `
		SELECT t.*
		FROM tasks t
		JOIN progress p ON t.id = p.task_id
		WHERE p.user_id = $1 AND t.level_id = $2 AND p.is_current = false
	`
	var tasks []domain.Task
	err := r.db.Select(&tasks, query, userID, levelID)
	return tasks, err
}

func (r *Repository) CountMistakes(userID, levelID int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM progress p
		JOIN tasks t ON p.task_id = t.id
		WHERE p.user_id = $1 AND t.level_id = $2 AND p.attempts > 1
	`
	var count int
	err := r.db.Get(&count, query, userID, levelID)
	return count, err
}

func (r *Repository) GetMatchPairsByTaskID(taskID int) ([]domain.TaskMatchPairs, error) {
	var pairs []domain.TaskMatchPairs

	query := `
		SELECT id, task_id, left_text, right_text, match_key
		FROM task_match_pairs
		WHERE task_id = $1
	`

	err := r.db.Select(&pairs, query, taskID)
	return pairs, err
}