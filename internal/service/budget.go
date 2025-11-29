package service

import (
	"context"
	"errors"
	"litespend-api/internal/model"
	"litespend-api/internal/repository"
	"time"
)

var (
	ErrBudgetNotFound = errors.New("budget not found")
)

type BudgetService struct {
	repo repository.BudgetRepository
}

func NewBudgetService(repository repository.BudgetRepository) *BudgetService {
	return &BudgetService{repo: repository}
}

func (s *BudgetService) Create(ctx context.Context, logined model.User, req model.CreateBudgetAllocationRequest) (int, error) {
	record := model.CreateBudgetAllocationRecord{
		UserID:     logined.ID,
		CategoryID: req.CategoryID,
		Year:       req.Year,
		Month:      req.Month,
		Assigned:   req.Assigned,
		CreatedAt:  time.Now(),
	}
	return s.repo.Create(ctx, record)
}

func (s *BudgetService) Update(ctx context.Context, logined model.User, id int, dto model.UpdateBudgetAllocationRequest) error {
	budget, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrBudgetNotFound
	}
	if budget.UserID != logined.ID && logined.Role != model.UserRoleAdmin {
		return ErrAccessDenied
	}
	return s.repo.Update(ctx, id, dto)
}

func (s *BudgetService) Delete(ctx context.Context, logined model.User, id int) error {
	budget, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrBudgetNotFound
	}
	if budget.UserID != logined.ID && logined.Role != model.UserRoleAdmin {
		return ErrAccessDenied
	}
	return s.repo.Delete(ctx, id)
}

func (s *BudgetService) GetByID(ctx context.Context, logined model.User, id int) (model.BudgetAllocation, error) {
	budget, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return budget, ErrBudgetNotFound
	}
	if budget.UserID != logined.ID && logined.Role != model.UserRoleAdmin {
		return budget, ErrAccessDenied
	}
	return budget, nil
}

func (s *BudgetService) GetList(ctx context.Context, logined model.User) ([]model.BudgetAllocation, error) {
	return s.repo.GetList(ctx, logined.ID)
}

func (s *BudgetService) GetListDetailedByPeriod(ctx context.Context, logined model.User, year uint, month uint) ([]model.BudgetDetailed, error) {
	return s.repo.GetListDetailedByPeriod(ctx, logined.ID, year, month)
}
