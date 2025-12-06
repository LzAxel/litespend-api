package repository

import (
	"context"
	"fmt"
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

func (r TransactionRepositoryPostgres) Create(ctx context.Context, transaction model.CreateTransactionRecord) (int, error) {
	var createdID int

	err := databases.WithinTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		err := tx.GetContext(ctx, &createdID,
			`INSERT INTO transactions(user_id, account_id, category_id, amount, date, note, approved, cleared, created_at, updated_at) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
			transaction.UserID,
			transaction.AccountID,
			transaction.CategoryID,
			transaction.Amount,
			transaction.Date,
			transaction.Note,
			transaction.IsApproved,
			transaction.IsCleared,
			transaction.CreatedAt,
			transaction.UpdatedAt,
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

func (r TransactionRepositoryPostgres) Update(ctx context.Context, id int, dto model.UpdateTransactionRecord) error {
	query := r.sq.Update("transactions").Where(sq.Eq{"id": id})

	if dto.AccountID != nil {
		query = query.Set("account_id", *dto.AccountID)
	}

	if dto.CategoryID != nil {
		query = query.Set("category_id", *dto.CategoryID)
	}

	if dto.Amount != nil {
		query = query.Set("amount", *dto.Amount)
	}

	if dto.Date != nil {
		query = query.Set("date", *dto.Date)
	}

	if dto.Note != nil {
		query = query.Set("note", *dto.Note)
	}

	if dto.IsCleared != nil {
		query = query.Set("cleared", *dto.IsCleared)
	}

	if dto.IsApproved != nil {
		query = query.Set("approved", *dto.IsApproved)
	}

	query.Set("updated_at", dto.UpdatedAt)

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
	var transactions []model.Transaction = make([]model.Transaction, 0)

	err := r.db.SelectContext(ctx, &transactions, `SELECT * FROM transactions WHERE user_id = $1 ORDER BY date DESC, created_at DESC`, userID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r TransactionRepositoryPostgres) GetListPaginated(ctx context.Context, userID uint64, accountID *uint64, params model.PaginationParams) ([]model.Transaction, int, error) {
	var transactions []model.Transaction = make([]model.Transaction, 0)
	var total int

	whereClause := "WHERE t.user_id = $1"
	args := []interface{}{userID}
	argIndex := 1

	if accountID != nil {
		argIndex++
		whereClause += fmt.Sprintf(" AND t.account_id = $%d", argIndex)
		args = append(args, *accountID)
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
