package databases

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"litespend-api/internal/config"
	"log/slog"
)

func GetPostgresPool(ctx context.Context, config config.PostgresConfig) (*pgxpool.Pool, error) {
	slog.InfoContext(ctx, "Connecting to Postgres")
	cfg, err := pgxpool.ParseConfig(config.GetDSNString())
	if err != nil {
		return nil, err
	}
	pgxPool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	if err = pgxPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return pgxPool, nil
}

func GetPostgresDB(pgxPool *pgxpool.Pool) *sqlx.DB {
	return sqlx.NewDb(stdlib.OpenDBFromPool(pgxPool), "pgx")
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
