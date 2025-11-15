package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"litespend-api/internal/model"
	"litespend-api/internal/repository/databases"
	"time"
)

type PrescribedExpanseRepositoryPostgres struct {
	db *sqlx.DB
	sq sq.StatementBuilderType
}

func NewPrescribedExpanseRepositoryPostgres(db *sqlx.DB) PrescribedExpanseRepositoryPostgres {
	return PrescribedExpanseRepositoryPostgres{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r PrescribedExpanseRepositoryPostgres) Create(ctx context.Context, prescribedExpanse model.CreatePrescribedExpanseRecord) (int, error) {
	var createdID int

	err := databases.WithinTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		err := tx.GetContext(ctx, &createdID,
			`INSERT INTO prescribed_expanses(user_id, category_id, description, frequency, amount, date_time, created_at) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			prescribedExpanse.UserID, prescribedExpanse.CategoryID, prescribedExpanse.Description,
			prescribedExpanse.Frequency, prescribedExpanse.Amount, prescribedExpanse.DateTime, prescribedExpanse.CreatedAt)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return createdID, nil
}

func (r PrescribedExpanseRepositoryPostgres) Update(ctx context.Context, id int, dto model.UpdatePrescribedExpanseRequest) error {
	query := r.sq.Update("prescribed_expanses").Where(sq.Eq{"id": id})

	if dto.CategoryID != nil {
		query = query.Set("category_id", *dto.CategoryID)
	}

	if dto.Description != nil {
		query = query.Set("description", *dto.Description)
	}

	if dto.Frequency != nil {
		query = query.Set("frequency", *dto.Frequency)
	}

	if dto.Amount != nil {
		query = query.Set("amount", *dto.Amount)
	}

	if dto.DateTime != nil {
		query = query.Set("date_time", *dto.DateTime)
	}

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r PrescribedExpanseRepositoryPostgres) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM prescribed_expanses WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}

func (r PrescribedExpanseRepositoryPostgres) GetByID(ctx context.Context, id int) (model.PrescribedExpanse, error) {
	var prescribedExpanse model.PrescribedExpanse

	err := r.db.GetContext(ctx, &prescribedExpanse, `SELECT * FROM prescribed_expanses WHERE id = $1`, id)
	if err != nil {
		return prescribedExpanse, err
	}

	return prescribedExpanse, nil
}

func (r PrescribedExpanseRepositoryPostgres) GetList(ctx context.Context, userID int) ([]model.PrescribedExpanse, error) {
	var prescribedExpanses []model.PrescribedExpanse

	err := r.db.SelectContext(ctx, &prescribedExpanses, `SELECT * FROM prescribed_expanses WHERE user_id = $1 ORDER BY date_time DESC`, userID)
	if err != nil {
		return prescribedExpanses, err
	}

	return prescribedExpanses, nil
}

func (r PrescribedExpanseRepositoryPostgres) IsPaidInPeriod(ctx context.Context, prescribedExpanseID int, periodStart, periodEnd time.Time) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count,
		`SELECT COUNT(*) FROM transactions 
		 WHERE prescribed_expanse_id = $1 
		 AND date_time >= $2 
		 AND date_time < $3`,
		prescribedExpanseID, periodStart, periodEnd)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r PrescribedExpanseRepositoryPostgres) GetPaidAmountInPeriod(ctx context.Context, prescribedExpanseID int, periodStart, periodEnd time.Time) (decimal.Decimal, *uint64, error) {
	// Сначала проверяем, есть ли транзакции
	var count int
	err := r.db.GetContext(ctx, &count,
		`SELECT COUNT(*) FROM transactions 
		 WHERE prescribed_expanse_id = $1 
		 AND date_time >= $2 
		 AND date_time < $3`,
		prescribedExpanseID, periodStart, periodEnd)
	if err != nil {
		return decimal.Zero, nil, err
	}

	if count == 0 {
		return decimal.Zero, nil, nil
	}

	// Если есть транзакции, получаем сумму всех транзакций и ID последней
	var totalAmount float64
	var lastTransactionID uint64

	// Получаем сумму всех транзакций
	err = r.db.GetContext(ctx, &totalAmount,
		`SELECT COALESCE(SUM(amount), 0)::float 
		 FROM transactions 
		 WHERE prescribed_expanse_id = $1 
		 AND date_time >= $2 
		 AND date_time < $3`,
		prescribedExpanseID, periodStart, periodEnd)
	if err != nil {
		return decimal.Zero, nil, err
	}

	// Получаем ID последней транзакции
	err = r.db.GetContext(ctx, &lastTransactionID,
		`SELECT id 
		 FROM transactions 
		 WHERE prescribed_expanse_id = $1 
		 AND date_time >= $2 
		 AND date_time < $3 
		 ORDER BY id DESC 
		 LIMIT 1`,
		prescribedExpanseID, periodStart, periodEnd)

	amount := decimal.NewFromFloat(totalAmount)
	var transactionID *uint64
	if err == nil && lastTransactionID > 0 {
		transactionID = &lastTransactionID
	}

	return amount, transactionID, nil
}
