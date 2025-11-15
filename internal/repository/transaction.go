package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"litespend-api/internal/model"
	"litespend-api/internal/repository/databases"
)

type TransactionRepositoryPostgres struct {
	db *sqlx.DB
	sq sq.StatementBuilderType
}

func NewTransactionRepositoryPostgres(db *sqlx.DB) TransactionRepositoryPostgres {
	return TransactionRepositoryPostgres{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r TransactionRepositoryPostgres) Create(ctx context.Context, transaction model.Transaction) (int, error) {
	var createdID int

	err := databases.WithinTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		err := tx.GetContext(ctx, &createdID,
			`INSERT INTO transactions(user_id, category_id, goal_id, description, amount, type, date_time, created_at) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
			transaction.UserID, transaction.CategoryID, transaction.GoalID, transaction.Description,
			transaction.Amount, transaction.Type, transaction.DateTime, transaction.CreatedAt)
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

func (r TransactionRepositoryPostgres) Update(ctx context.Context, id int, dto model.UpdateTransactionRequest) error {
	query := r.sq.Update("transactions").Where(sq.Eq{"id": id})

	if dto.CategoryID != nil {
		query = query.Set("category_id", *dto.CategoryID)
	}

	if dto.GoalID != nil {
		query = query.Set("goal_id", *dto.GoalID)
	}

	if dto.Description != nil {
		query = query.Set("description", *dto.Description)
	}

	if dto.Amount != nil {
		query = query.Set("amount", *dto.Amount)
	}

	if dto.Type != nil {
		query = query.Set("type", *dto.Type)
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

func (r TransactionRepositoryPostgres) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM transactions WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}

func (r TransactionRepositoryPostgres) GetByID(ctx context.Context, id int) (model.Transaction, error) {
	var transaction model.Transaction

	err := r.db.GetContext(ctx, &transaction, `SELECT * FROM transactions WHERE id = $1`, id)
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r TransactionRepositoryPostgres) GetList(ctx context.Context, userID int) ([]model.Transaction, error) {
	var transactions []model.Transaction

	err := r.db.SelectContext(ctx, &transactions, `SELECT * FROM transactions WHERE user_id = $1 ORDER BY date_time DESC`, userID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}
