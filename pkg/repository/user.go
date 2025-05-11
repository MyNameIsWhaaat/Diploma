package repository

import "github.com/MyNameIsWhaaat/algo-learning/pkg/domain"

func (r *Repository) GetUserProfileData(userID int) (domain.UserProfileData, error) {
	query := `
		SELECT 
			u.id, u.username, u.email, u.avatar_url,
			pl.title AS profile_level,
			upl.total_xp
		FROM users u
		JOIN user_profile_levels upl ON u.id = upl.user_id
		JOIN profile_levels pl ON upl.profile_level_id = pl.id
		WHERE u.id = $1
	`
	var data domain.UserProfileData
	err := r.db.Get(&data, query, userID)
	return data, err
}