package repository

import (
	"context"

	"litespend-api/internal/model"
	"litespend-api/internal/repository/databases"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
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

func (r BudgetRepositoryPostgres) Create(ctx context.Context, record model.CreateBudgetAllocationRecord) (int, error) {
	var createdID int

	err := databases.WithinTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		err := tx.GetContext(ctx, &createdID,
			`INSERT INTO budget_allocations(user_id, category_id, year, month, assigned, created_at) 
			 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
			record.UserID, record.CategoryID, record.Year, record.Month, record.Assigned, record.CreatedAt,
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

func (r BudgetRepositoryPostgres) Update(ctx context.Context, id int, dto model.UpdateBudgetAllocationRequest) error {
	query := r.sq.Update("budget_allocations").Where(sq.Eq{"id": id})

	if dto.CategoryID != nil {
		query = query.Set("category_id", *dto.CategoryID)
	}
	if dto.Year != nil {
		query = query.Set("year", *dto.Year)
	}
	if dto.Month != nil {
		query = query.Set("month", *dto.Month)
	}
	if dto.Assigned != nil {
		query = query.Set("assigned", *dto.Assigned)
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
	_, err := r.db.ExecContext(ctx, `DELETE FROM budget_allocations WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (r BudgetRepositoryPostgres) GetByID(ctx context.Context, id int) (model.BudgetAllocation, error) {
	var b model.BudgetAllocation
	err := r.db.GetContext(ctx, &b, `SELECT * FROM budget_allocations WHERE id = $1`, id)
	if err != nil {
		return b, err
	}
	return b, nil
}

func (r BudgetRepositoryPostgres) GetList(ctx context.Context, userID uint64) ([]model.BudgetAllocation, error) {
	var items []model.BudgetAllocation = make([]model.BudgetAllocation, 0)
	err := r.db.SelectContext(ctx, &items, `SELECT * FROM budget_allocations WHERE user_id = $1 ORDER BY year DESC, month DESC, category_id`, userID)
	if err != nil {
		return items, err
	}
	return items, nil
}

func (r BudgetRepositoryPostgres) GetListDetailedByPeriod(ctx context.Context, userID uint64, year uint64, month uint64) (model.CategoryBudgetResponse, error) {
	query := `
		WITH target AS (
			SELECT
				$1::int AS target_year,
				$2::int AS target_month,
				$3::bigint AS user_id
		),
		month_bounds AS (
			SELECT
				make_date(t.target_year, t.target_month, 1)                                      AS month_start,
				(make_date(t.target_year, t.target_month, 1) + interval '1 month' - interval '1 day') AS month_end
			FROM target t
		),
		prev_month AS (
			SELECT
				CASE WHEN t.target_month = 1 THEN t.target_year - 1 ELSE t.target_year END      AS prev_year,
				CASE WHEN t.target_month = 1 THEN 12 ELSE t.target_month - 1 END               AS prev_month
			FROM target t
		),
		prev_bounds AS (
			SELECT
				make_date(pm.prev_year, pm.prev_month, 1)                                      AS prev_start,
				(make_date(pm.prev_year, pm.prev_month, 1) + interval '1 month' - interval '1 day') AS prev_end
			FROM prev_month pm
		),
		account_balances AS (
			SELECT COALESCE(SUM(t.amount), 0)::numeric AS total_balance
			FROM transactions t
			WHERE t.user_id = (SELECT user_id FROM target)
		),
		total_assigned AS (
			SELECT COALESCE(SUM(ba.assigned), 0)::numeric AS assigned_sum
			FROM budget_allocations ba, target t
			WHERE ba.user_id = t.user_id
				AND ba.year = t.target_year
				AND ba.month = t.target_month
		),
		category_data AS (
			SELECT
				c.id                                      AS category_id,
				c.name                                    AS category_name,
				COALESCE(c.group_name, '')                AS category_group_name,
				COALESCE(ba.assigned, 0)::numeric         AS assigned,

				-- Потрачено в текущем месяце
				COALESCE((
					SELECT SUM(t.amount) * -1
					FROM transactions t
					CROSS JOIN month_bounds mb
					WHERE t.user_id = tar.user_id
						AND t.category_id = c.id
						AND t.amount < 0
						AND t.date BETWEEN mb.month_start AND mb.month_end
				), 0)::numeric AS spent,

				-- Перенос с прошлого месяца (assigned - spent)
				COALESCE((
					SELECT COALESCE(SUM(prev_ba.assigned), 0) - COALESCE(SUM(prev_t.amount * -1), 0)
					FROM budget_allocations prev_ba
					LEFT JOIN transactions prev_t
						ON prev_t.category_id = c.id
						AND prev_t.user_id = tar.user_id
						AND prev_t.amount < 0
						AND prev_t.date BETWEEN pb.prev_start AND pb.prev_end
					CROSS JOIN prev_bounds pb
					WHERE prev_ba.user_id = tar.user_id
						AND prev_ba.category_id = c.id
						AND prev_ba.year = pm.prev_year
						AND prev_ba.month = pm.prev_month
				), 0)::numeric AS carried_over

			FROM categories c
			CROSS JOIN target tar
			CROSS JOIN prev_month pm
			CROSS JOIN prev_bounds pb
			LEFT JOIN budget_allocations ba
				ON ba.category_id = c.id
				AND ba.user_id = tar.user_id
				AND ba.year = tar.target_year
				AND ba.month = tar.target_month
			WHERE c.user_id = tar.user_id
		)
		SELECT
			(ab.total_balance - ta.assigned_sum)::numeric AS tbb,
			cd.category_id,
			cd.category_name,
			cd.category_group_name,
			cd.assigned,
			cd.spent,
			(cd.assigned + cd.carried_over - cd.spent)::numeric AS available,
			cd.carried_over
		FROM category_data cd
		CROSS JOIN account_balances ab
		CROSS JOIN total_assigned ta
		ORDER BY cd.category_name;
	`

	rows, err := r.db.QueryxContext(ctx, query, year, month, userID)
	if err != nil {
		return model.CategoryBudgetResponse{}, err
	}
	defer rows.Close()

	response := model.CategoryBudgetResponse{}

	first := true
	for rows.Next() {
		var cat model.CategoryBudget
		var tbb decimal.Decimal

		if first {
			err = rows.Scan(
				&tbb,
				&cat.CategoryID,
				&cat.Name,
				&cat.GroupName,
				&cat.Assigned,
				&cat.Spent,
				&cat.Available,
				&cat.CarriedOver,
			)
			first = false
		} else {
			err = rows.Scan(
				new(float64), // пропускаем tbb
				&cat.CategoryID,
				&cat.Name,
				&cat.GroupName,
				&cat.Assigned,
				&cat.Spent,
				&cat.Available,
				&cat.CarriedOver,
			)
		}
		if err != nil {
			return response, err
		}

		response.ToBeBudgeted = tbb
		response.Categories = append(response.Categories, cat)
	}
	return response, nil
}
