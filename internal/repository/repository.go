package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"litespend-api/internal/model"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user model.CreateUserRecord) (int, error)
	Update(ctx context.Context, id int, dto model.UpdateUserRecord) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (model.User, error)
	GetByUsername(ctx context.Context, username string) (model.User, error)
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction model.Transaction) (int, error)
	Update(ctx context.Context, id int, dto model.UpdateTransactionRequest) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (model.Transaction, error)
	GetList(ctx context.Context, userID uint64) ([]model.Transaction, error)
	GetListPaginated(ctx context.Context, userID uint64, params model.PaginationParams) ([]model.Transaction, int, error)
	GetBalanceStatistics(ctx context.Context, userID uint64) (model.CurrentBalanceStatistics, error)
	GetCategoryStatistics(ctx context.Context, userID uint64, period model.PeriodType, from, to *time.Time) ([]model.CategoryStatisticsItem, error)
	GetPeriodStatistics(ctx context.Context, userID uint64, period model.PeriodType, from, to *time.Time) ([]model.PeriodStatisticsItem, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, category model.CreateCategoryRecord) (int, error)
	Update(ctx context.Context, id int, dto model.UpdateCategoryRequest) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (model.TransactionCategory, error)
	GetList(ctx context.Context, userID uint64) ([]model.TransactionCategory, error)
	GetListByType(ctx context.Context, userID uint64, categoryType model.CategoryType) ([]model.TransactionCategory, error)
}

type BudgetRepository interface {
	Create(ctx context.Context, record model.CreateBudgetRecord) (int, error)
	Update(ctx context.Context, id int, dto model.UpdateBudgetRequest) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (model.Budget, error)
	GetList(ctx context.Context, userID uint64) ([]model.Budget, error)
	GetListByPeriod(ctx context.Context, userID uint64, year uint, month uint) ([]model.Budget, error)
}

type Repository struct {
	UserRepository        UserRepository
	TransactionRepository TransactionRepository
	CategoryRepository    CategoryRepository
	BudgetRepository      BudgetRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepository:        NewUserRepositoryPostgres(db),
		TransactionRepository: NewTransactionRepositoryPostgres(db),
		CategoryRepository:    NewCategoryRepositoryPostgres(db),
		BudgetRepository:      NewBudgetRepositoryPostgres(db),
	}
}
