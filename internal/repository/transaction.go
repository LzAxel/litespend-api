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

func (r TransactionRepositoryPostgres) GetListPaginated(ctx context.Context, userID int, params model.PaginationParams) ([]model.Transaction, int, error) {
	var transactions []model.Transaction
	var total int

	err := r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM transactions WHERE user_id = $1`, userID)
	if err != nil {
		return transactions, 0, err
	}

	err = r.db.SelectContext(ctx, &transactions,
		`SELECT * FROM transactions WHERE user_id = $1 ORDER BY date_time DESC LIMIT $2 OFFSET $3`,
		userID, params.Limit, params.Offset())
	if err != nil {
		return transactions, 0, err
	}

	return transactions, total, nil
}

func (r TransactionRepositoryPostgres) GetBalanceStatistics(ctx context.Context, userID int) (model.CurrentBalanceStatistics, error) {
	var stats model.CurrentBalanceStatistics
	var totalIncome, totalExpense *float64

	row := r.db.QueryRowContext(ctx,
		`SELECT 
			COALESCE(SUM(CASE WHEN type = $1 THEN amount ELSE 0 END), 0) as total_income,
			COALESCE(SUM(CASE WHEN type = $2 THEN amount ELSE 0 END), 0) as total_expense
		FROM transactions 
		WHERE user_id = $3`,
		model.TransactionTypeIncome, model.TransactionTypeExpanse, userID)

	err := row.Scan(&totalIncome, &totalExpense)
	if err != nil {
		return stats, err
	}

	income := float64(0)
	expense := float64(0)
	if totalIncome != nil {
		income = *totalIncome
	}
	if totalExpense != nil {
		expense = *totalExpense
	}

	stats.TotalIncome = decimal.NewFromFloat(income)
	stats.TotalExpense = decimal.NewFromFloat(expense)
	stats.Balance = stats.TotalIncome.Sub(stats.TotalExpense)

	return stats, nil
}

func (r TransactionRepositoryPostgres) GetCategoryStatistics(ctx context.Context, userID int, period model.PeriodType, from, to *time.Time) ([]model.CategoryStatisticsItem, error) {
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
			t.category_id,
			tc.name as category_name,
			TO_CHAR(date_trunc('%s', t.date_time), '%s') as period,
			COALESCE(SUM(CASE WHEN t.type = $1 THEN t.amount ELSE 0 END), 0)::float as income,
			COALESCE(SUM(CASE WHEN t.type = $2 THEN t.amount ELSE 0 END), 0)::float as expense
		FROM transactions t
		INNER JOIN transaction_categories tc ON t.category_id = tc.id
		WHERE t.user_id = $3
	`, periodFormat, periodCharFormat)

	args := []interface{}{
		model.TransactionTypeIncome,
		model.TransactionTypeExpanse,
		userID,
	}
	argIndex := 3

	if from != nil {
		argIndex++
		query += fmt.Sprintf(" AND t.date_time >= $%d", argIndex)
		args = append(args, *from)
	}

	if to != nil {
		argIndex++
		query += fmt.Sprintf(" AND t.date_time <= $%d", argIndex)
		args = append(args, *to)
	}

	query += fmt.Sprintf(`
		GROUP BY t.category_id, tc.name, date_trunc('%s', t.date_time)
		ORDER BY period DESC, tc.name
	`, periodFormat)

	err := r.db.SelectContext(ctx, &rows, query, args...)
	if err != nil {
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

func (r TransactionRepositoryPostgres) GetPeriodStatistics(ctx context.Context, userID int, period model.PeriodType, from, to *time.Time) ([]model.PeriodStatisticsItem, error) {
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
			TO_CHAR(date_trunc('%s', t.date_time), '%s') as period,
			COALESCE(SUM(CASE WHEN t.type = $1 THEN t.amount ELSE 0 END), 0)::float as income,
			COALESCE(SUM(CASE WHEN t.type = $2 THEN t.amount ELSE 0 END), 0)::float as expense
		FROM transactions t
		WHERE t.user_id = $3
	`, periodFormat, periodCharFormat)

	args := []interface{}{
		model.TransactionTypeIncome,
		model.TransactionTypeExpanse,
		userID,
	}
	argIndex := 3

	if from != nil {
		argIndex++
		query += fmt.Sprintf(" AND t.date_time >= $%d", argIndex)
		args = append(args, *from)
	}

	if to != nil {
		argIndex++
		query += fmt.Sprintf(" AND t.date_time <= $%d", argIndex)
		args = append(args, *to)
	}

	query += fmt.Sprintf(`
		GROUP BY date_trunc('%s', t.date_time)
		ORDER BY period DESC
	`, periodFormat)

	err := r.db.SelectContext(ctx, &rows, query, args...)
	if err != nil {
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
