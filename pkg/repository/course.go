package repository

import (

	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
)

func (r *Repository) GetCourseWithProgress(userId int) ([]domain.CourseWithProgress, error) {
	query := `
		SELECT
			c.id,
			c.title,
			c.short_description,
			c.image_url,
			uc.xp_earned,
			c.xp_reward,
			uc.completed
		FROM user_courses uc
		JOIN courses c ON uc.course_id = c.id
		WHERE uc.user_id = $1
		ORDER BY c.id;
	`

	var courses []domain.CourseWithProgress
	err := r.db.Select(&courses, query, userId)
	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (r *Repository) GetCourseWithoutProgress(userId int) ([]domain.CourseWithoutProgress, error) {
	query := `
		SELECT
			id,
			title,
			short_description,
			full_description,
			image_url,
			xp_reward
		FROM courses
		WHERE NOT EXISTS (
			SELECT 1 FROM user_courses
			WHERE user_courses.course_id = courses.id AND user_courses.user_id = $1
		)
		ORDER BY id;
	`

	var courses []domain.CourseWithoutProgress
	err := r.db.Select(&courses, query, userId)
	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (r *Repository) StartCourse(userID, courseID int) error {
	_, err := r.db.Exec(`
	INSERT INTO user_courses (user_id, course_id, xp_earned, completed)
	VALUES ($1, $2, 0, false)
	ON CONFLICT (user_id, course_id) DO NOTHING
`, userID, courseID)

	return err
}
