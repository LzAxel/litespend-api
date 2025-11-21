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

func (r TransactionRepositoryPostgres) GetBalanceStatistics(ctx context.Context, userID uint64) (model.CurrentBalanceStatistics, error) {
	var stats model.CurrentBalanceStatistics

	query := `
WITH vars AS (
    SELECT $1::bigint AS uid,
           EXTRACT(YEAR FROM CURRENT_DATE)::int  AS y,
           EXTRACT(MONTH FROM CURRENT_DATE)::int AS m
),
on_accounts AS (
    SELECT COALESCE(SUM(amount), 0)::numeric(14,2) AS balance
    FROM transactions WHERE user_id = (SELECT uid FROM vars)
),
unpaid_bills AS (
    SELECT COALESCE(SUM(amount_expected - amount_paid), 0)::numeric(14,2) AS reserved
    FROM bill_instances bi
    JOIN recurring_bills rb ON bi.recurring_bill_id = rb.id
    WHERE rb.user_id = (SELECT uid FROM vars)
      AND rb.is_active
      AND (bi.year > (SELECT y FROM vars) OR (bi.year = (SELECT y FROM vars) AND bi.month >= (SELECT m FROM vars)))
),
budgeted_this_month AS (
    SELECT COALESCE(SUM(budgeted), 0)::numeric(14,2) AS reserved
    FROM monthly_budgets
    WHERE user_id = (SELECT uid FROM vars)
      AND year = (SELECT y FROM vars)
      AND month = (SELECT m FROM vars)
)
SELECT 
    on_accounts.balance AS on_accounts,
    unpaid_bills.reserved AS reserved_bills,
    budgeted_this_month.reserved AS reserved_budgets
FROM vars, on_accounts, unpaid_bills, budgeted_this_month`

	type row struct {
		OnAccounts     decimal.Decimal `db:"on_accounts"`
		ReservedBills  decimal.Decimal `db:"reserved_bills"`
		ReservedBudget decimal.Decimal `db:"reserved_budgets"`
	}

	var snapshot row
	if err := r.db.GetContext(ctx, &snapshot, query, userID); err != nil {
		return stats, err
	}

	stats.OnAccounts = snapshot.OnAccounts
	stats.ReservedBills = snapshot.ReservedBills
	stats.ReservedBudgets = snapshot.ReservedBudget
	stats.TotalReserved = snapshot.ReservedBills.Add(snapshot.ReservedBudget)
	stats.FreeToDistribute = snapshot.OnAccounts.Sub(stats.TotalReserved)

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
