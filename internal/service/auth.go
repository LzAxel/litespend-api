package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"litespend-api/internal/model"
	"litespend-api/internal/repository"
	"litespend-api/internal/session"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrForbidden       = errors.New("forbidden")
)

type AuthService struct {
	sessionManager *session.SessionManager
	userRepo       repository.UserRepository
}

func NewAuthService(sessionManager *session.SessionManager, userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		sessionManager: sessionManager,
		userRepo:       userRepo,
	}
}

func (s *AuthService) RevokeSession(ctx context.Context, logined model.User, token string) error {
	if logined.Role != model.UserRoleAdmin {
		return ErrForbidden
	}

	exists, err := s.sessionManager.SessionExists(token)
	if err != nil {
		return err
	}

	if !exists {
		return ErrSessionNotFound
	}

	return s.sessionManager.DeleteByToken(token)
}

func (s *AuthService) GetSessionInfo(ctx context.Context, logined model.User, token string) (model.SessionInfo, error) {
	if logined.Role != model.UserRoleAdmin {
		return model.SessionInfo{}, ErrForbidden
	}

	exists, err := s.sessionManager.SessionExists(token)
	if err != nil {
		return model.SessionInfo{}, err
	}

	if !exists {
		return model.SessionInfo{}, ErrSessionNotFound
	}

	return model.SessionInfo{
		Token: token,
	}, nil
}
