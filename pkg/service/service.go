package service

import (
	"github.com/MyNameIsWhaaat/algo-learning"
	"github.com/MyNameIsWhaaat/algo-learning/pkg/repository"
)

type Authorization interface {
	CreateUser(user algolearning.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}