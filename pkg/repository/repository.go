package repository

import (
	"github.com/MyNameIsWhaaat/algo-learning"
	"github.com/jmoiron/sqlx"
)

type Authorization interface{
	CreateUser(user algolearning.User) (int, error)
	GetUser(username, password string) (algolearning.User, error)
}

type Repository struct{
	Authorization
	
}

func NewRepository(db *sqlx.DB)  *Repository{
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}