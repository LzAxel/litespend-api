package service

import (
	"context"
	"litespend-api/internal/model"
	"litespend-api/internal/repository"
)

type Service struct {
	User
}

type User interface {
	Register(ctx context.Context, user model.RegisterRequest) error
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		User: NewUserService(repository.UserRepository),
	}
}
