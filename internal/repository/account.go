package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"litespend-api/internal/model"
)

type AccountRepositoryPostgres struct {
	db *sqlx.DB
	sq sq.StatementBuilderType
}

func NewAccountRepositoryPostgres(db *sqlx.DB) AccountRepositoryPostgres {
	return AccountRepositoryPostgres{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r AccountRepositoryPostgres) Create(ctx context.Context, account model.CreateAccountRecord) (uint64, error) {
	var createdID uint64

	err := r.db.GetContext(ctx, &createdID, `
			INSERT INTO accounts (user_id, name, type, is_archived, order_num, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5, $6, $7) 
			RETURNING id`,
		account.UserID, account.Name, account.Type, account.IsArchived, account.OrderNum, account.CreatedAt, account.UpdatedAt,
	)

	if err != nil {
		return 0, err
	}

	return createdID, nil
}

func (r AccountRepositoryPostgres) Update(ctx context.Context, id uint64, dto model.UpdateAccountRecord) error {
	query := r.sq.Update("accounts").Where(sq.Eq{"id": id})

	if dto.Name != nil {
		query = query.Set("name", *dto.Name)
	}

	if dto.IsArchived != nil {
		query = query.Set("is_archived", *dto.IsArchived)
	}

	if dto.OrderNum != nil {
		query = query.Set("order_num", *dto.OrderNum)
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

func (r AccountRepositoryPostgres) Delete(ctx context.Context, id uint64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM accounts WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}

func (r AccountRepositoryPostgres) GetByID(ctx context.Context, id uint64) (model.AccountDB, error) {
	var account model.AccountDB

	err := r.db.GetContext(ctx, &account, `
		SELECT a.*, SUM(tr.amount) as balance FROM accounts a
		         LEFT JOIN transactions tr ON a.id = tr.account_id
		WHERE a.id = $1`, id)
	if err != nil {
		return account, err
	}

	return account, nil
}

func (r AccountRepositoryPostgres) GetList(ctx context.Context, userID uint64) ([]model.AccountDB, error) {
	var accounts []model.AccountDB = make([]model.AccountDB, 0)

	err := r.db.SelectContext(ctx, &accounts, `
		SELECT a.*, SUM(tr.amount) as balance FROM accounts a 
		         LEFT JOIN transactions tr ON a.id = tr.account_id
		WHERE a.user_id = $1 
		ORDER BY a.order_num, a.name`, userID)
	if err != nil {
		return accounts, err
	}

	return accounts, nil
}
