package service

import (
	"context"
	"litespend-api/internal/model"
	"litespend-api/internal/repository"
	"time"
)

type AccountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repository repository.AccountRepository) *AccountService {
	return &AccountService{repo: repository}
}

func (s *AccountService) Create(ctx context.Context, logined model.User, account model.CreateAccountRequest) (uint64, error) {
	createdID, err := s.repo.Create(ctx, model.CreateAccountRecord{
		UserID:     logined.ID,
		Name:       account.Name,
		Type:       account.Type,
		IsArchived: account.IsArchived,
		OrderNum:   account.OrderNum,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})
	if err != nil {
		return 0, err
	}

	return createdID, nil
}

func (s *AccountService) GetList(ctx context.Context, logined model.User) ([]model.Account, error) {
	accounts, err := s.repo.GetList(ctx, logined.ID)
	if err != nil {
		return []model.Account{}, err
	}

	return accounts, nil
}

func (s *AccountService) Delete(ctx context.Context, logined model.User, id uint64) error {
	account, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if account.UserID != logined.ID && logined.Role != model.UserRoleAdmin {
		return ErrAccessDenied
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *AccountService) Update(ctx context.Context, logined model.User, id uint64, dto model.UpdateAccountRequest) error {
	account, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if account.UserID != logined.ID && logined.Role != model.UserRoleAdmin {
		return ErrAccessDenied
	}

	err = s.repo.Update(ctx, id, model.UpdateAccountRecord{
		Name:       dto.Name,
		IsArchived: dto.IsArchived,
		OrderNum:   dto.OrderNum,
		UpdatedAt:  time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}
