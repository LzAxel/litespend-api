package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"litespend-api/internal/model"
	"litespend-api/internal/repository/databases"
)

type UserRepositoryPostgres struct {
	db *sqlx.DB
	sq sq.StatementBuilderType
}

func NewUserRepositoryPostgres(db *sqlx.DB) UserRepositoryPostgres {
	return UserRepositoryPostgres{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r UserRepositoryPostgres) Create(ctx context.Context, user model.CreateUserRecord) (int, error) {
	var createdID int

	err := databases.WithinTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		err := tx.GetContext(ctx, &createdID, `INSERT INTO users(username, role, password_hash, created_at) VALUES ($1, $2, $3, $4) RETURNING id`, user.Username, user.Role, user.PasswordHash, user.CreatedAt)
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

func (r UserRepositoryPostgres) Update(ctx context.Context, id int, dto model.UpdateUserRecord) error {
	query := r.sq.Update("users").Where(sq.Eq{"id": id})

	if dto.Username != nil {
		query = query.Set("username", *dto.Username)
	}

	if dto.Role != nil {
		query = query.Set("role", *dto.Role)
	}

	if dto.PasswordHash != nil {
		query = query.Set("password_hash", *dto.PasswordHash)
	}

	sqlQuery, args, _ := query.ToSql()

	_, err := r.db.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepositoryPostgres) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepositoryPostgres) GetByID(ctx context.Context, id int) (model.User, error) {
	var user model.User

	err := r.db.GetContext(ctx, &user, `SELECT * FROM users WHERE id = $1`, id)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r UserRepositoryPostgres) GetByUsername(ctx context.Context, username string) (model.User, error) {
	var user model.User

	err := r.db.GetContext(ctx, &user, `SELECT * FROM users WHERE username = $1`, username)
	if err != nil {
		return user, err
	}

	return user, nil
}
