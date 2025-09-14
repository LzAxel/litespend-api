package service

import (
	"context"
	"litespend-api/internal/model"
	"litespend-api/internal/repository"
	"litespend-api/internal/session"
)

type Service struct {
	User
}

type User interface {
	Register(ctx context.Context, user model.RegisterRequest) error
	Login(ctx context.Context, user model.LoginRequest) error
}

func NewService(repository *repository.Repository, sm *session.SessionManager) *Service {
	return &Service{
		User: NewUserService(repository.UserRepository, sm),
	}
}
