package service

import (
	"context"
	"errors"
	"time"

	"litespend-api/internal/model"
	"litespend-api/internal/repository"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrAccessDenied        = errors.New("access denied")
)

type TransactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repository repository.TransactionRepository) *TransactionService {
	return &TransactionService{
		repo: repository,
	}
}

func (s *TransactionService) Create(ctx context.Context, logined model.User, req model.CreateTransactionRequest) (int, error) {
	transaction := model.CreateTransactionRecord{
		UserID:     logined.ID,
		CategoryID: req.CategoryID,
		Amount:     req.Amount,
		Date:       req.Date,
		AccountID:  req.AccountID,
		Note:       req.Note,
		IsCleared:  req.IsCleared,
		IsApproved: req.IsApproved,
		UpdatedAt:  time.Now(),
		CreatedAt:  time.Now(),
	}

	id, err := s.repo.Create(ctx, transaction)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *TransactionService) Update(ctx context.Context, logined model.User, id int, dto model.UpdateTransactionRequest) error {
	transaction, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrTransactionNotFound
	}

	if transaction.UserID != logined.ID && logined.Role != model.UserRoleAdmin {
		return ErrAccessDenied
	}

	err = s.repo.Update(ctx, id, model.UpdateTransactionRecord{
		AccountID:  dto.AccountID,
		CategoryID: dto.CategoryID,
		Amount:     dto.Amount,
		Date:       dto.Date,
		Note:       dto.Note,
		IsCleared:  dto.IsCleared,
		IsApproved: dto.IsApproved,
		UpdatedAt:  time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *TransactionService) Delete(ctx context.Context, logined model.User, id int) error {
	transaction, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrTransactionNotFound
	}

	if transaction.UserID != logined.ID && logined.Role != model.UserRoleAdmin {
		return ErrAccessDenied
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *TransactionService) GetByID(ctx context.Context, logined model.User, id int) (model.Transaction, error) {
	transaction, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return transaction, ErrTransactionNotFound
	}

	if transaction.UserID != logined.ID && logined.Role != model.UserRoleAdmin {
		return transaction, ErrAccessDenied
	}

	return transaction, nil
}

func (s *TransactionService) GetList(ctx context.Context, logined model.User) ([]model.Transaction, error) {
	transactions, err := s.repo.GetList(ctx, logined.ID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *TransactionService) GetListPaginated(ctx context.Context, logined model.User, accountID *uint64, params model.PaginationParams) (model.PaginatedTransactionsResponse, error) {
	params.Validate()

	transactions, total, err := s.repo.GetListPaginated(ctx, logined.ID, accountID, params)
	if err != nil {
		return model.PaginatedTransactionsResponse{}, err
	}

	return model.NewPaginatedResponse(transactions, total, params), nil
}
