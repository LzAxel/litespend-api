package service

import (
	"context"
	"errors"
	"litespend-api/internal/model"
	"litespend-api/internal/repository"
	"time"
)

var (
	ErrPrescribedExpanseNotFound = errors.New("prescribed expanse not found")
)

type PrescribedExpanseService struct {
	repo repository.PrescribedExpanseRepository
}

func NewPrescribedExpanseService(repository repository.PrescribedExpanseRepository) *PrescribedExpanseService {
	return &PrescribedExpanseService{
		repo: repository,
	}
}

func (s *PrescribedExpanseService) Create(ctx context.Context, logined model.User, req model.CreatePrescribedExpanseRequest) (int, error) {
	prescribedExpanse := model.CreatePrescribedExpanseRecord{
		UserID:      logined.ID,
		CategoryID:  req.CategoryID,
		Description: req.Description,
		Frequency:   req.Frequency,
		Amount:      req.Amount,
		DateTime:    req.DateTime,
		CreatedAt:   time.Now(),
	}

	id, err := s.repo.Create(ctx, prescribedExpanse)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *PrescribedExpanseService) Update(ctx context.Context, logined model.User, id int, dto model.UpdatePrescribedExpanseRequest) error {
	prescribedExpanse, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrPrescribedExpanseNotFound
	}

	if prescribedExpanse.UserID != logined.ID && logined.Role != model.UserRoleAdmin {
		return ErrAccessDenied
	}

	err = s.repo.Update(ctx, id, dto)
	if err != nil {
		return err
	}

	return nil
}

func (s *PrescribedExpanseService) Delete(ctx context.Context, logined model.User, id int) error {
	prescribedExpanse, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrPrescribedExpanseNotFound
	}

	if prescribedExpanse.UserID != logined.ID && logined.Role != model.UserRoleAdmin {
		return ErrAccessDenied
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *PrescribedExpanseService) GetByID(ctx context.Context, logined model.User, id int) (model.PrescribedExpanse, error) {
	prescribedExpanse, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return prescribedExpanse, ErrPrescribedExpanseNotFound
	}

	if prescribedExpanse.UserID != logined.ID && logined.Role != model.UserRoleAdmin {
		return prescribedExpanse, ErrAccessDenied
	}

	return prescribedExpanse, nil
}

func (s *PrescribedExpanseService) GetList(ctx context.Context, logined model.User) ([]model.PrescribedExpanse, error) {
	prescribedExpanses, err := s.repo.GetList(ctx, int(logined.ID))
	if err != nil {
		return prescribedExpanses, err
	}

	return prescribedExpanses, nil
}
