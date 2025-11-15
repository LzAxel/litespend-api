package service

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
	"litespend-api/internal/model"
	"litespend-api/internal/repository"
	"time"
)

var (
	ErrPrescribedExpanseNotFound = errors.New("prescribed expanse not found")
)

type PrescribedExpanseService struct {
	repo            repository.PrescribedExpanseRepository
	transactionRepo repository.TransactionRepository
}

func NewPrescribedExpanseService(repository repository.PrescribedExpanseRepository, transactionRepo repository.TransactionRepository) *PrescribedExpanseService {
	return &PrescribedExpanseService{
		repo:            repository,
		transactionRepo: transactionRepo,
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

func (s *PrescribedExpanseService) GetListWithPaymentStatus(ctx context.Context, logined model.User) ([]model.PrescribedExpanseWithPaymentStatus, error) {
	prescribedExpanses, err := s.repo.GetList(ctx, int(logined.ID))
	if err != nil {
		return nil, err
	}

	now := time.Now()
	result := make([]model.PrescribedExpanseWithPaymentStatus, len(prescribedExpanses))
	for i, pe := range prescribedExpanses {
		periodStart, periodEnd := s.getPeriodForFrequency(pe.Frequency, pe.DateTime, now)
		paidAmount, transactionID, err := s.repo.GetPaidAmountInPeriod(ctx, int(pe.ID), periodStart, periodEnd)
		if err != nil {
			return nil, err
		}
		// Трата считается оплаченной, если сумма оплаты >= сумме обязательной траты
		isPaid := paidAmount.GreaterThanOrEqual(pe.Amount)
		result[i] = model.PrescribedExpanseWithPaymentStatus{
			PrescribedExpanse: pe,
			IsPaid:            isPaid,
			PaidAmount:        paidAmount,
			TransactionID:     transactionID,
		}
	}

	return result, nil
}

func (s *PrescribedExpanseService) MarkAsPaid(ctx context.Context, logined model.User, id int) (int, error) {
	prescribedExpanse, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return 0, ErrPrescribedExpanseNotFound
	}

	if prescribedExpanse.UserID != logined.ID && logined.Role != model.UserRoleAdmin {
		return 0, ErrAccessDenied
	}

	now := time.Now()
	prescribedExpanseID := uint64(id)
	transaction := model.Transaction{
		UserID:              logined.ID,
		CategoryID:          prescribedExpanse.CategoryID,
		PrescribedExpanseID: &prescribedExpanseID,
		Description:         prescribedExpanse.Description,
		Amount:              prescribedExpanse.Amount,
		Type:                model.TransactionTypeExpanse,
		DateTime:            now,
		CreatedAt:           now,
	}

	transactionID, err := s.transactionRepo.Create(ctx, transaction)
	if err != nil {
		return 0, err
	}

	return transactionID, nil
}

func (s *PrescribedExpanseService) MarkAsPaidPartial(ctx context.Context, logined model.User, id int, amount decimal.Decimal) (int, error) {
	prescribedExpanse, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return 0, ErrPrescribedExpanseNotFound
	}

	if prescribedExpanse.UserID != logined.ID && logined.Role != model.UserRoleAdmin {
		return 0, ErrAccessDenied
	}

	if amount.LessThanOrEqual(decimal.Zero) {
		return 0, errors.New("сумма должна быть больше нуля")
	}

	if amount.GreaterThan(prescribedExpanse.Amount) {
		return 0, errors.New("сумма не может быть больше суммы обязательной траты")
	}

	now := time.Now()
	prescribedExpanseID := uint64(id)
	transaction := model.Transaction{
		UserID:              logined.ID,
		CategoryID:          prescribedExpanse.CategoryID,
		PrescribedExpanseID: &prescribedExpanseID,
		Description:         prescribedExpanse.Description + " (частичная оплата)",
		Amount:              amount,
		Type:                model.TransactionTypeExpanse,
		DateTime:            now,
		CreatedAt:           now,
	}

	transactionID, err := s.transactionRepo.Create(ctx, transaction)
	if err != nil {
		return 0, err
	}

	return transactionID, nil
}

func (s *PrescribedExpanseService) getPeriodForFrequency(frequency model.FrequencyType, startDate time.Time, currentDate time.Time) (time.Time, time.Time) {
	var periodStart, periodEnd time.Time

	switch frequency {
	case model.FrequencyTypeDaily:
		periodStart = time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 0, 0, 0, 0, currentDate.Location())
		periodEnd = periodStart.AddDate(0, 0, 1)
	case model.FrequencyTypeWeekly:
		weekday := int(currentDate.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		periodStart = currentDate.AddDate(0, 0, -weekday+1)
		periodStart = time.Date(periodStart.Year(), periodStart.Month(), periodStart.Day(), 0, 0, 0, 0, periodStart.Location())
		periodEnd = periodStart.AddDate(0, 0, 7)
	case model.FrequencyTypeMonthly:
		periodStart = time.Date(currentDate.Year(), currentDate.Month(), 1, 0, 0, 0, 0, currentDate.Location())
		periodEnd = periodStart.AddDate(0, 1, 0)
	case model.FrequencyTypeQuarterly:
		quarter := (int(currentDate.Month()) - 1) / 3
		periodStart = time.Date(currentDate.Year(), time.Month(quarter*3+1), 1, 0, 0, 0, 0, currentDate.Location())
		periodEnd = periodStart.AddDate(0, 3, 0)
	default:
		periodStart = time.Date(currentDate.Year(), currentDate.Month(), 1, 0, 0, 0, 0, currentDate.Location())
		periodEnd = periodStart.AddDate(0, 1, 0)
	}

	return periodStart, periodEnd
}
