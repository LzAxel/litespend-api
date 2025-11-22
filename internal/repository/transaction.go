package repository

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"litespend-api/internal/model"
	"litespend-api/internal/repository/databases"
	"time"
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
			`INSERT INTO transactions(user_id, category_id, amount, date, description, created_at, bill_instance_id) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			transaction.UserID, transaction.CategoryID, transaction.Amount, transaction.Date, transaction.Description, transaction.CreatedAt, transaction.BillInstanceID)
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

	if dto.Description != nil {
		query = query.Set("description", *dto.Description)
	}

	if dto.Amount != nil {
		query = query.Set("amount", *dto.Amount)
	}

	if dto.Date != nil {
		query = query.Set("date", *dto.Date)
	}

	if dto.BillInstanceID != nil {
		query = query.Set("bill_instance_id", *dto.BillInstanceID)
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

func (r TransactionRepositoryPostgres) GetList(ctx context.Context, userID uint64) ([]model.Transaction, error) {
	var transactions []model.Transaction

	err := r.db.SelectContext(ctx, &transactions, `SELECT * FROM transactions WHERE user_id = $1 ORDER BY date DESC, created_at DESC`, userID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r TransactionRepositoryPostgres) GetListPaginated(ctx context.Context, userID uint64, params model.PaginationParams) ([]model.Transaction, int, error) {
	var transactions []model.Transaction
	var total int

	whereClause := "WHERE t.user_id = $1"
	args := []interface{}{userID}
	argIndex := 1

	if params.Search != nil && *params.Search != "" {
		argIndex++
		whereClause += fmt.Sprintf(" AND t.description ILIKE $%d", argIndex)
		args = append(args, "%"+*params.Search+"%")
	}

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM transactions t %s`, whereClause)
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return transactions, 0, err
	}

	orderBy := "ORDER BY t.date DESC, t.created_at DESC"
	if params.SortBy != nil {
		sortOrder := "DESC"
		if params.SortOrder != nil && *params.SortOrder == model.SortOrderASC {
			sortOrder = "ASC"
		}

		switch *params.SortBy {
		case model.SortFieldDate:
			orderBy = fmt.Sprintf("ORDER BY t.date %s", sortOrder)
		case model.SortFieldDescription:
			orderBy = fmt.Sprintf("ORDER BY t.description %s", sortOrder)
		case model.SortFieldCategory:
			orderBy = fmt.Sprintf("ORDER BY t.category_id %s", sortOrder)
		}
	}

	argIndex++
	limitArg := argIndex
	argIndex++
	offsetArg := argIndex
	args = append(args, params.Limit, params.Offset())

	query := fmt.Sprintf(`
		SELECT t.* FROM transactions t 
		%s 
		%s 
		LIMIT $%d OFFSET $%d
	`, whereClause, orderBy, limitArg, offsetArg)

	err = r.db.SelectContext(ctx, &transactions, query, args...)
	if err != nil {
		return transactions, 0, err
	}

	return transactions, total, nil
}

func (r TransactionRepositoryPostgres) GetBalanceStatistics(ctx context.Context, userID uint64, year uint, month uint) (model.CurrentBalanceStatistics, error) {
	var stats model.CurrentBalanceStatistics

	query := `
		WITH 
		balance AS (
			SELECT COALESCE(SUM(CASE WHEN amount < 0 THEN amount ELSE 0 END), 0)::numeric(14,2) as expanses,
			COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0)::numeric(14,2) as income
			FROM transactions WHERE user_id = $1
		),
		budgets AS (
			SELECT
				COALESCE(SUM(b.budgeted), 0)::numeric(14,2) AS reserved,
				COALESCE(SUM(t.amount), 0)::numeric(14,2) AS spent,
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
		)
		SELECT 
			balance.expanses AS total_expanse,
			balance.income AS total_income,
			budgets.reserved AS total_reserved,
			budgets.reserved AS reserved_budgets,
			balance.income - balance.expanses - budgets.reserved AS free_to_distribute
		FROM budgets, balance`

	type row struct {
		TotalExpense     decimal.Decimal `db:"total_expanse"`
		TotalIncome      decimal.Decimal `db:"total_income"`
		TotalReserved    decimal.Decimal `db:"total_reserved"`
		ReservedBudgets  decimal.Decimal `db:"reserved_budgets"`
		FreeToDistribute decimal.Decimal `db:"free_to_distribute"`
	}

	var snapshot row
	if err := r.db.GetContext(ctx, &snapshot, query, userID, year, month); err != nil {
		return stats, err
	}

	stats.TotalExpense = snapshot.TotalExpense
	stats.TotalIncome = snapshot.TotalIncome
	stats.TotalReserved = snapshot.TotalReserved
	stats.ReservedBudgets = snapshot.ReservedBudgets
	stats.FreeToDistribute = snapshot.FreeToDistribute

	return stats, nil
}

func (r TransactionRepositoryPostgres) GetCategoryStatistics(ctx context.Context, userID uint64, period model.PeriodType, from, to *time.Time) ([]model.CategoryStatisticsItem, error) {
	type categoryStatsRow struct {
		CategoryID   uint64   `db:"category_id"`
		CategoryName string   `db:"category_name"`
		Period       string   `db:"period"`
		Income       *float64 `db:"income"`
		Expense      *float64 `db:"expense"`
	}

	var rows []categoryStatsRow
	var periodFormat, periodCharFormat string

	switch period {
	case model.PeriodTypeDay:
		periodFormat = "day"
		periodCharFormat = "YYYY-MM-DD"
	case model.PeriodTypeWeek:
		periodFormat = "week"
		periodCharFormat = "YYYY-MM-DD"
	case model.PeriodTypeMonth:
		periodFormat = "month"
		periodCharFormat = "YYYY-MM"
	default:
		periodFormat = "day"
		periodCharFormat = "YYYY-MM-DD"
	}

	query := fmt.Sprintf(`
		SELECT 
			COALESCE(t.category_id, 0) as category_id,
			COALESCE(c.name, 'Без категории') as category_name,
			TO_CHAR(date_trunc('%s', t.date), '%s') as period,
			COALESCE(SUM(CASE WHEN t.amount > 0 THEN t.amount ELSE 0 END), 0)::float as income,
			COALESCE(SUM(CASE WHEN t.amount < 0 THEN ABS(t.amount) ELSE 0 END), 0)::float as expense
		FROM transactions t
		LEFT JOIN transaction_categories c ON t.category_id = c.id
		WHERE t.user_id = $1
	`, periodFormat, periodCharFormat)

	args := []interface{}{userID}
	argIndex := 1

	if from != nil {
		argIndex++
		query += fmt.Sprintf(" AND t.date >= $%d", argIndex)
		args = append(args, *from)
	}

	if to != nil {
		argIndex++
		query += fmt.Sprintf(" AND t.date <= $%d", argIndex)
		args = append(args, *to)
	}

	query += fmt.Sprintf(`
		GROUP BY COALESCE(t.category_id, 0), c.name, date_trunc('%s', t.date)
		ORDER BY period DESC, category_name
	`, periodFormat)

	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	items := make([]model.CategoryStatisticsItem, len(rows))
	for i, row := range rows {
		incomeVal := float64(0)
		expenseVal := float64(0)
		if row.Income != nil {
			incomeVal = *row.Income
		}
		if row.Expense != nil {
			expenseVal = *row.Expense
		}

		items[i] = model.CategoryStatisticsItem{
			CategoryID:   row.CategoryID,
			CategoryName: row.CategoryName,
			Period:       row.Period,
			Income:       decimal.NewFromFloat(incomeVal),
			Expense:      decimal.NewFromFloat(expenseVal),
		}
	}

	return items, nil
}

func (r TransactionRepositoryPostgres) GetPeriodStatistics(ctx context.Context, userID uint64, period model.PeriodType, from, to *time.Time) ([]model.PeriodStatisticsItem, error) {
	type periodStatsRow struct {
		Period  string   `db:"period"`
		Income  *float64 `db:"income"`
		Expense *float64 `db:"expense"`
	}

	var rows []periodStatsRow
	var periodFormat, periodCharFormat string

	switch period {
	case model.PeriodTypeDay:
		periodFormat = "day"
		periodCharFormat = "YYYY-MM-DD"
	case model.PeriodTypeWeek:
		periodFormat = "week"
		periodCharFormat = "YYYY-MM-DD"
	case model.PeriodTypeMonth:
		periodFormat = "month"
		periodCharFormat = "YYYY-MM"
	default:
		periodFormat = "day"
		periodCharFormat = "YYYY-MM-DD"
	}

	query := fmt.Sprintf(`
		SELECT 
			TO_CHAR(date_trunc('%s', t.date), '%s') as period,
			COALESCE(SUM(CASE WHEN t.amount > 0 THEN t.amount ELSE 0 END), 0)::float as income,
			COALESCE(SUM(CASE WHEN t.amount < 0 THEN ABS(t.amount) ELSE 0 END), 0)::float as expense
		FROM transactions t
		WHERE t.user_id = $1
	`, periodFormat, periodCharFormat)

	args := []interface{}{userID}
	argIndex := 1

	if from != nil {
		argIndex++
		query += fmt.Sprintf(" AND t.date >= $%d", argIndex)
		args = append(args, *from)
	}

	if to != nil {
		argIndex++
		query += fmt.Sprintf(" AND t.date <= $%d", argIndex)
		args = append(args, *to)
	}

	query += fmt.Sprintf(`
		GROUP BY date_trunc('%s', t.date)
		ORDER BY period DESC
	`, periodFormat)

	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, err
	}

	items := make([]model.PeriodStatisticsItem, len(rows))
	for i, row := range rows {
		incomeVal := float64(0)
		expenseVal := float64(0)
		if row.Income != nil {
			incomeVal = *row.Income
		}
		if row.Expense != nil {
			expenseVal = *row.Expense
		}

		income := decimal.NewFromFloat(incomeVal)
		expense := decimal.NewFromFloat(expenseVal)
		balance := income.Sub(expense)

		items[i] = model.PeriodStatisticsItem{
			Period:  row.Period,
			Income:  income,
			Expense: expense,
			Balance: balance,
		}
	}

	return items, nil
}
