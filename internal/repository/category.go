package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"litespend-api/internal/model"
	"litespend-api/internal/repository/databases"
)

type CategoryRepositoryPostgres struct {
	db *sqlx.DB
	sq sq.StatementBuilderType
}

func NewCategoryRepositoryPostgres(db *sqlx.DB) CategoryRepositoryPostgres {
	return CategoryRepositoryPostgres{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r CategoryRepositoryPostgres) Create(ctx context.Context, category model.CreateCategoryRecord) (int, error) {
	var createdID int

	err := databases.WithinTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		err := tx.GetContext(ctx, &createdID, `INSERT INTO transaction_categories(user_id, name, type, created_at) VALUES ($1, $2, $3, $4) RETURNING id`, category.UserID, category.Name, category.Type, category.CreatedAt)
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

func (r CategoryRepositoryPostgres) Update(ctx context.Context, id int, dto model.UpdateCategoryRequest) error {
	query := r.sq.Update("transaction_categories").Where(sq.Eq{"id": id})

	if dto.Name != nil {
		query = query.Set("name", *dto.Name)
	}

	if dto.Type != nil {
		query = query.Set("type", *dto.Type)
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

func (r CategoryRepositoryPostgres) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM transaction_categories WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}

func (r CategoryRepositoryPostgres) GetByID(ctx context.Context, id int) (model.TransactionCategory, error) {
	var category model.TransactionCategory

	err := r.db.GetContext(ctx, &category, `SELECT * FROM transaction_categories WHERE id = $1`, id)
	if err != nil {
		return category, err
	}

	return category, nil
}

func (r CategoryRepositoryPostgres) GetList(ctx context.Context, userID uint64) ([]model.TransactionCategory, error) {
	var categories []model.TransactionCategory = make([]model.TransactionCategory, 0)

	err := r.db.SelectContext(ctx, &categories, `SELECT * FROM transaction_categories WHERE user_id = $1 ORDER BY name`, userID)
	if err != nil {
		return categories, err
	}

	return categories, nil
}

func (r CategoryRepositoryPostgres) GetListByType(ctx context.Context, userID uint64, categoryType model.CategoryType) ([]model.TransactionCategory, error) {
	var categories []model.TransactionCategory = make([]model.TransactionCategory, 0)

	err := r.db.SelectContext(ctx, &categories, `SELECT * FROM transaction_categories WHERE user_id = $1 AND type = $2 ORDER BY name`, userID, categoryType)
	if err != nil {
		return categories, err
	}

	return categories, nil
}
