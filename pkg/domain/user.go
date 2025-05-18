package domain

import "time"

type User struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Username  string    `db:"username"`
	Password  string    `db:"password_hash"`
	Email     *string   `db:"email"`
	AvatarURL *string   `db:"avatar_url"`
	Role      string    `db:"role"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	IsNewUser bool      `db:"is_new_user"`
}
