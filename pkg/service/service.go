package service

import (
	"github.com/MyNameIsWhaaat/algo-learning/pkg/domain"
)

//type LayerService interface {
//	CreateUser(user domain.User) (int, error)
//	GenerateToken(username, password string) (string, error)
//	ParseToken(token string) (int, error)
//
//	List() ([]domain.Course, error)
//}

type Repository interface {
	CreateUser(user domain.User) (int, error)
	GetUser(username, password string) (domain.User, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}
