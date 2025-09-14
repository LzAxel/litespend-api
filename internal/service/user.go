package service

import (
	"context"
	"litespend-api/internal/model"
	"litespend-api/internal/pkg/hash"
	"litespend-api/internal/repository"
	"litespend-api/internal/session"
	"log/slog"
	"time"
)

type UserService struct {
	repo           repository.UserRepository
	sessionManager *session.SessionManager
}

func NewUserService(repository repository.UserRepository, sm *session.SessionManager) *UserService {
	return &UserService{
		repo:           repository,
		sessionManager: sm,
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

func (s *UserService) Login(ctx context.Context, user model.LoginRequest) error {
	foundUser, err := s.repo.GetByUsername(ctx, user.Username)
	if err != nil {
		return err
	}

	if err := hash.CheckPassword(foundUser.PasswordHash, user.Password); err != nil {
		return err
	}
	slog.Info("token", s.sessionManager.GetManager().Token(ctx))

	return nil
}
