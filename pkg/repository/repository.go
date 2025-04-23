package repository

import (
	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
	"github.com/jmoiron/sqlx"
)

type Authorization interface{
	CreateUser(user domain.User) (int, error)
	GetUser(username, password string) (domain.User, error)
}

type Repository struct{
	Authorization
	
}

func NewRepository(db *sqlx.DB)  *Repository{
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}