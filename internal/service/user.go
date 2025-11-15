package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"litespend-api/internal/model"
	"litespend-api/internal/pkg/hash"
	"litespend-api/internal/repository"
	"litespend-api/internal/session"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

type UserService struct {
	repo           repository.UserRepository
	sessionManager *session.SessionManager
}

func NewUserService(repository repository.UserRepository, sessionManager *session.SessionManager) *UserService {
	return &UserService{
		repo:           repository,
		sessionManager: sessionManager,
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

func (s *UserService) Login(ctx context.Context, c *gin.Context, req model.LoginRequest) (model.User, error) {
	user, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return model.User{}, ErrInvalidCredentials
	}

	err = hash.CheckPassword(user.PasswordHash, req.Password)
	if err != nil {
		return model.User{}, ErrInvalidCredentials
	}

	if s.sessionManager != nil {
		err = s.sessionManager.RenewToken(c)
		if err != nil {
			return model.User{}, err
		}

		s.sessionManager.Put(c, "user_id", int(user.ID))
	}

	return user, nil
}

func (s *UserService) Logout(ctx context.Context, c *gin.Context) error {
	if s.sessionManager != nil {
		return s.sessionManager.Destroy(c)
	}
	return nil
}
