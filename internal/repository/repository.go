package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"litespend-api/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user model.CreateUserRecord) (int, error)
	Update(ctx context.Context, id int, dto model.UpdateUserRecord) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (model.User, error)
	GetByUsername(ctx context.Context, username string) (model.User, error)
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction model.CreateTransactionRecord) (int, error)
	Replace(ctx context.Context, id int, dto model.UpdateTransactionRecord) error
	Delete(ctx context.Context, id int) error
	GetListByUserID(ctx context.Context, userID int) ([]model.Transaction, error)
	GetByID(ctx context.Context, id int) (model.Transaction, error)
}

type TransactionCategoryRepository interface {
	Create(ctx context.Context, category model.CreateTransactionCategoryRecord) (int, error)
	Replace(ctx context.Context, id int, dto model.UpdateTransactionCategoryRecord) error
	Delete(ctx context.Context, id int) error
	GetListByUserID(ctx context.Context, userID int) ([]model.TransactionCategory, error)
	GetByID(ctx context.Context, id int) (model.TransactionCategory, error)
}

type GoalRepository interface {
	Create(ctx context.Context, goal model.CreateGoalRecord) (int, error)
	Replace(ctx context.Context, id int, dto model.UpdateGoalRecord) error
	Delete(ctx context.Context, id int) error
	GetListByUserID(ctx context.Context, userID int) ([]model.Goal, error)
	GetByID(ctx context.Context, id int) (model.Goal, error)
}

type Repository struct {
	UserRepository
	TransactionCategoryRepository
	TransactionRepository
	GoalRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepository:                NewUserRepositoryPostgres(db),
		TransactionCategoryRepository: NewTransactionCategoryRepositoryPostgres(db),
		TransactionRepository:         NewTransactionRepositoryPostgres(db),
		GoalRepository:                NewGoalRepositoryPostgres(db),
	}
}
