package service

import (
	"context"
	"litespend-api/internal/model"
	"litespend-api/internal/pkg/hash"
	"litespend-api/internal/repository"
	"litespend-api/internal/session"
	"time"
)

type UserService struct {
	repo           repository.UserRepository
	sessionManager *session.SessionManager
}

func NewUserService(repository repository.UserRepository) *UserService {
	return &UserService{
		repo: repository,
	}
}

func (s *UserService) Register(ctx context.Context, user model.RegisterRequest) error {
	hashedPassword, err := hash.HashPassword(user.Password)
	if err != nil {
		return err
	}

	_, err = s.repo.Create(ctx, model.CreateUserRecord{
		Username:     user.Username,
		Role:         model.UserRoleUser,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}
