package databases

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"litespend-api/internal/config"
	"log/slog"
)

func GetPostgresDB(ctx context.Context, config config.PostgresConfig) (*sqlx.DB, error) {
	slog.InfoContext(ctx, "Connecting to Postgres")
	db, err := sqlx.ConnectContext(ctx, "pgx", config.GetDSNString())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return db, nil
}

func WithinTransaction(ctx context.Context, db *sqlx.DB, fn func(tx *sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func CheckSqZeroSetStatementsError(err error) bool {
	return err.Error() == "update statements must have at least one Set clause"
}
