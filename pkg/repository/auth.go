package repository

import (
	"fmt"

	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
)

// CreateUser создает пользователя, возвращает его идентификатор.
func (r *Repository) CreateUser(user domain.User) (int, error) {
	const query = `INSERT INTO users (name, username, email, avatar_url, password_hash, is_new_user) values ($1, $2, $3, $4, $5, $6) RETURNING id`

	row := r.db.QueryRow(query, user.Name, user.Username, user.Email, "0", user.Password, true)

	var userID int
	if err := row.Scan(&userID); err != nil {
		return 0, fmt.Errorf("error creating user: row.Scan: %w", err)
	}

	//Создаем запись в таблице с уровнями пользователей
	profileLevelID := 1
	_, err := r.db.Exec(`
			INSERT INTO user_profile_levels (user_id, profile_level_id, total_xp, last_level_up)
			VALUES ($1, $2, 0, NOW())
		`, userID, profileLevelID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *Repository) GetUser(username, password string) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", "users")
	err := r.db.Get(&user, query, username, password)

	return user, err
}
