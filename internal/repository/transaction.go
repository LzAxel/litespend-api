package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"litespend-api/internal/model"
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

func (r TransactionRepositoryPostgres) Create(ctx context.Context, transaction model.CreateTransactionRecord) (int, error) {
	var createdID int
	sql := `INSERT INTO transactions(
			user_id,    
			category_id,
			goal_id,    
			description,
			amount,     
			type,       
			date_time,  
    		created_at
		) VALUE ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	err := r.db.GetContext(ctx, &createdID, sql,
		transaction.UserID,
		transaction.CategoryID,
		transaction.GoalID,
		transaction.Description,
		transaction.Amount,
		transaction.Type,
		transaction.DateTime,
		transaction.CreatedAt)
	if err != nil {
		return 0, err
	}

	return createdID, nil
}

func (r TransactionRepositoryPostgres) Replace(ctx context.Context, id int, dto model.UpdateTransactionRecord) error {
	sql := `UPDATE transactions(category_id, goal_id, description, amount, type, date_time) WHERE id=$1 VALUES ($2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, sql, id, dto.CategoryID, dto.GoalID, dto.Description, dto.Amount, dto.Type, dto.DateTime)
	return err
}

func (r TransactionRepositoryPostgres) Delete(ctx context.Context, id int) error {
	sql := `DELETE FROM transactions WHERE id=$1`

	_, err := r.db.ExecContext(ctx, sql, id)
	return err
}

func (r TransactionRepositoryPostgres) GetListByUserID(ctx context.Context, userID int) ([]model.Transaction, error) {
	var transactions = make([]model.Transaction, 0)

	sql := `SELECT * FROM transactions WHERE user_id=$1`

	err := r.db.SelectContext(ctx, transactions, sql, userID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r TransactionRepositoryPostgres) GetByID(ctx context.Context, id int) (model.Transaction, error) {
	var transaction model.Transaction

	sql := `SELECT * FROM transactions WHERE id=$1`

	err := r.db.GetContext(ctx, &transaction, sql, id)
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}
