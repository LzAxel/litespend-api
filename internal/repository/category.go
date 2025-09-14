package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"litespend-api/internal/model"
)

type TransactionCategoryRepositoryPostgres struct {
	db *sqlx.DB
	sq sq.StatementBuilderType
}

func NewTransactionCategoryRepositoryPostgres(db *sqlx.DB) TransactionCategoryRepositoryPostgres {
	return TransactionCategoryRepositoryPostgres{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r TransactionCategoryRepositoryPostgres) Create(ctx context.Context, category model.CreateTransactionCategoryRecord) (int, error) {
	var createdID int
	sql := `INSERT INTO transaction_categories(name, user_id, type, created_at) VALUE ($1, $2, $3, $4) RETURNING id`

	err := r.db.GetContext(ctx, &createdID, sql, category.Name, category.UserID, category.Type, category.CreatedAt)
	if err != nil {
		return 0, err
	}

	return createdID, nil
}

func (r TransactionCategoryRepositoryPostgres) Replace(ctx context.Context, id int, dto model.UpdateTransactionCategoryRecord) error {
	sql := `UPDATE transaction_categories(name, type) WHERE id=$1 VALUES ($2, $3)`

	_, err := r.db.ExecContext(ctx, sql, id, dto.Name, dto.Type)
	return err
}

func (r TransactionCategoryRepositoryPostgres) Delete(ctx context.Context, id int) error {
	sql := `DELETE FROM transaction_categories WHERE id=$1`

	_, err := r.db.ExecContext(ctx, sql, id)
	return err
}

func (r TransactionCategoryRepositoryPostgres) GetListByUserID(ctx context.Context, userID int) ([]model.TransactionCategory, error) {
	var categories = make([]model.TransactionCategory, 0)

	sql := `SELECT * FROM transaction_categories WHERE user_id=$1`

	err := r.db.SelectContext(ctx, categories, sql, userID)
	if err != nil {
		return categories, err
	}

	return categories, nil
}

func (r TransactionCategoryRepositoryPostgres) GetByID(ctx context.Context, id int) (model.TransactionCategory, error) {
	var category model.TransactionCategory

	sql := `SELECT * FROM transaction_categories WHERE id=$1`

	err := r.db.GetContext(ctx, &category, sql, id)
	if err != nil {
		return category, err
	}

	return category, nil
}
