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
		err := tx.GetContext(ctx, &createdID, `INSERT INTO categories(user_id, name, group_name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`, category.UserID, category.Name, category.GroupName, category.CreatedAt, category.UpdatedAt)
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

func (r CategoryRepositoryPostgres) Update(ctx context.Context, id int, dto model.UpdateCategoryRecord) error {
	query := r.sq.Update("categories").Where(sq.Eq{"id": id})

	if dto.Name != nil {
		query = query.Set("name", *dto.Name)
	}

	if dto.GroupName != nil {
		query = query.Set("group_name", *dto.GroupName)
	}

	query = query.Set("updated_at", dto.UpdatedAt)

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
	_, err := r.db.ExecContext(ctx, `DELETE FROM categories WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}

func (r CategoryRepositoryPostgres) GetByID(ctx context.Context, id int) (model.Category, error) {
	var category model.Category

	err := r.db.GetContext(ctx, &category, `SELECT * FROM categories WHERE id = $1`, id)
	if err != nil {
		return category, err
	}

	return category, nil
}

func (r CategoryRepositoryPostgres) GetList(ctx context.Context, userID uint64) ([]model.Category, error) {
	var categories []model.Category = make([]model.Category, 0)

	err := r.db.SelectContext(ctx, &categories, `SELECT * FROM categories WHERE user_id = $1 ORDER BY name`, userID)
	if err != nil {
		return categories, err
	}

	return categories, nil
}
