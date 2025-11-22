package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"litespend-api/internal/model"
	"litespend-api/internal/repository/databases"
)

type BudgetRepositoryPostgres struct {
	db *sqlx.DB
	sq sq.StatementBuilderType
}

func NewBudgetRepositoryPostgres(db *sqlx.DB) BudgetRepositoryPostgres {
	return BudgetRepositoryPostgres{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r BudgetRepositoryPostgres) Create(ctx context.Context, record model.CreateBudgetRecord) (int, error) {
	var createdID int

	err := databases.WithinTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		err := tx.GetContext(ctx, &createdID,
			`INSERT INTO budgets(user_id, category_id, year, month, budgeted, created_at) 
			 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
			record.UserID, record.CategoryID, record.Year, record.Month, record.Budgeted, record.CreatedAt,
		)
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

func (r BudgetRepositoryPostgres) Update(ctx context.Context, id int, dto model.UpdateBudgetRequest) error {
	query := r.sq.Update("budgets").Where(sq.Eq{"id": id})

	if dto.CategoryID != nil {
		query = query.Set("category_id", *dto.CategoryID)
	}
	if dto.Year != nil {
		query = query.Set("year", *dto.Year)
	}
	if dto.Month != nil {
		query = query.Set("month", *dto.Month)
	}
	if dto.Budgeted != nil {
		query = query.Set("budgeted", *dto.Budgeted)
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

func (r BudgetRepositoryPostgres) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM budgets WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (r BudgetRepositoryPostgres) GetByID(ctx context.Context, id int) (model.Budget, error) {
	var b model.Budget
	err := r.db.GetContext(ctx, &b, `SELECT * FROM budgets WHERE id = $1`, id)
	if err != nil {
		return b, err
	}
	return b, nil
}

func (r BudgetRepositoryPostgres) GetList(ctx context.Context, userID uint64) ([]model.Budget, error) {
	var items []model.Budget = make([]model.Budget, 0)
	err := r.db.SelectContext(ctx, &items, `SELECT * FROM budgets WHERE user_id = $1 ORDER BY year DESC, month DESC, category_id`, userID)
	if err != nil {
		return items, err
	}
	return items, nil
}

func (r BudgetRepositoryPostgres) GetListByPeriod(ctx context.Context, userID uint64, year uint, month uint) ([]model.Budget, error) {
	var items []model.Budget = make([]model.Budget, 0)
	err := r.db.SelectContext(ctx, &items, `SELECT * FROM budgets WHERE user_id = $1 AND year = $2 AND month = $3 ORDER BY category_id`, userID, year, month)
	if err != nil {
		return items, err
	}
	return items, nil
}

func (r BudgetRepositoryPostgres) GetListDetailedByPeriod(ctx context.Context, userID uint64, year uint, month uint) ([]model.BudgetDetailed, error) {
	var items = make([]model.BudgetDetailed, 0)
	err := r.db.SelectContext(ctx, &items,
		`SELECT 
    		b.*,
            COALESCE(SUM(t.amount), 0)::numeric::float8 AS spent,
            (b.budgeted + COALESCE(SUM(t.amount), 0))::numeric::float8 AS remaining
        FROM budgets b
        LEFT JOIN transaction_categories c ON c.id = b.category_id
        LEFT JOIN transactions t 
            ON t.category_id = b.category_id 
            AND t.user_id = b.user_id
            AND DATE_TRUNC('month', t.date) = make_date($2, $3, 1)
        WHERE 
            b.user_id = $1
            AND b.year = $2
            AND b.month = $3
        GROUP BY 
            b.id, b.category_id, c.name, b.budgeted
        ORDER BY 
            c.name`,
		userID, year, month)
	if err != nil {
		return items, err
	}
	return items, nil
}
