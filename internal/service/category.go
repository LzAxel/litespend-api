package service

import (
	"context"
	"errors"
	"litespend-api/internal/model"
	"litespend-api/internal/repository"
	"time"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
)

type CategoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repository repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		repo: repository,
	}
}

func (s *CategoryService) Create(ctx context.Context, logined model.User, req model.CreateCategoryRequest) (int, error) {
	category := model.CreateCategoryRecord{
		UserID:    logined.ID,
		Name:      req.Name,
		Type:      req.Type,
		CreatedAt: time.Now(),
	}

	id, err := s.repo.Create(ctx, category)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *CategoryService) Update(ctx context.Context, logined model.User, id int, dto model.UpdateCategoryRequest) error {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrCategoryNotFound
	}

	if category.UserID != logined.ID && logined.Role != model.UserRoleAdmin {
		return ErrAccessDenied
	}

	err = s.repo.Update(ctx, id, dto)
	if err != nil {
		return err
	}

	return nil
}

func (s *CategoryService) Delete(ctx context.Context, logined model.User, id int) error {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrCategoryNotFound
	}

	if category.UserID != logined.ID && logined.Role != model.UserRoleAdmin {
		return ErrAccessDenied
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *CategoryService) GetByID(ctx context.Context, logined model.User, id int) (model.TransactionCategory, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return category, ErrCategoryNotFound
	}

	if category.UserID != logined.ID && logined.Role != model.UserRoleAdmin {
		return category, ErrAccessDenied
	}

	return category, nil
}

func (s *CategoryService) GetList(ctx context.Context, logined model.User) ([]model.TransactionCategory, error) {
	categories, err := s.repo.GetList(ctx, logined.ID)
	if err != nil {
		return categories, err
	}

	return categories, nil
}

func (s *CategoryService) GetListByType(ctx context.Context, logined model.User, categoryType model.CategoryType) ([]model.TransactionCategory, error) {
	categories, err := s.repo.GetListByType(ctx, logined.ID, categoryType)
	if err != nil {
		return categories, err
	}

	return categories, nil
}
