package service

import (
	"context"
	"errors"
	"litespend-api/internal/model"
	"litespend-api/internal/pkg/hash"
	"litespend-api/internal/repository"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

type UserService struct {
	repo repository.UserRepository
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

func (s *UserService) Login(ctx context.Context, req model.LoginRequest) (model.User, error) {
	user, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return model.User{}, ErrInvalidCredentials
	}

	err = hash.CheckPassword(user.PasswordHash, req.Password)
	if err != nil {
		return model.User{}, ErrInvalidCredentials
	}

	return user, nil
}

func (s *UserService) Logout(ctx context.Context) error {
	return nil
}
