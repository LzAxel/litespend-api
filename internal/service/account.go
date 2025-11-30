package service

import "litespend-api/internal/repository"

type AccountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repository repository.AccountRepository) *AccountService {
	return &AccountService{repo: repository}
}
