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
	Update(ctx context.Context, id int, dto model.UpdateTransactionRecord) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (model.Transaction, error)
	GetList(ctx context.Context, userID uint64) ([]model.Transaction, error)
	GetListPaginated(ctx context.Context, userID uint64, accountID *uint64, params model.PaginationParams) ([]model.Transaction, int, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, category model.CreateCategoryRecord) (int, error)
	Update(ctx context.Context, id int, dto model.UpdateCategoryRecord) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (model.Category, error)
	GetList(ctx context.Context, userID uint64) ([]model.Category, error)
}

type BudgetRepository interface {
	Create(ctx context.Context, record model.CreateBudgetAllocationRecord) (int, error)
	Update(ctx context.Context, id int, dto model.UpdateBudgetAllocationRequest) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (model.BudgetAllocation, error)
	GetList(ctx context.Context, userID uint64) ([]model.BudgetAllocation, error)
	GetListByPeriod(ctx context.Context, userID uint64, year uint, month uint) ([]model.BudgetAllocation, error)
}

type AccountRepository interface {
	Create(ctx context.Context, account model.CreateAccountRecord) (int, error)
	Update(ctx context.Context, id int, dto model.UpdateAccountRecord) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (model.Account, error)
	GetList(ctx context.Context, userID uint64) ([]model.Account, error)
}

type Repository struct {
	UserRepository        UserRepository
	TransactionRepository TransactionRepository
	CategoryRepository    CategoryRepository
	BudgetRepository      BudgetRepository
	AccountRepository     AccountRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepository:        NewUserRepositoryPostgres(db),
		TransactionRepository: NewTransactionRepositoryPostgres(db),
		CategoryRepository:    NewCategoryRepositoryPostgres(db),
		BudgetRepository:      NewBudgetRepositoryPostgres(db),
		AccountRepository:     NewAccountRepositoryPostgres(db),
	}
}
