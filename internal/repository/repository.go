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
	Create(ctx context.Context, transaction model.Transaction) (int, error)
	Update(ctx context.Context, id int, dto model.UpdateTransactionRequest) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (model.Transaction, error)
	GetList(ctx context.Context, userID int) ([]model.Transaction, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, category model.CreateCategoryRecord) (int, error)
	Update(ctx context.Context, id int, dto model.UpdateCategoryRequest) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (model.TransactionCategory, error)
	GetList(ctx context.Context, userID int) ([]model.TransactionCategory, error)
	GetListByType(ctx context.Context, userID int, categoryType model.TransactionType) ([]model.TransactionCategory, error)
}

type PrescribedExpanseRepository interface {
	Create(ctx context.Context, prescribedExpanse model.CreatePrescribedExpanseRecord) (int, error)
	Update(ctx context.Context, id int, dto model.UpdatePrescribedExpanseRequest) error
	Delete(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (model.PrescribedExpanse, error)
	GetList(ctx context.Context, userID int) ([]model.PrescribedExpanse, error)
}

type Repository struct {
	UserRepository              UserRepository
	TransactionRepository       TransactionRepository
	CategoryRepository          CategoryRepository
	PrescribedExpanseRepository PrescribedExpanseRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepository:              NewUserRepositoryPostgres(db),
		TransactionRepository:       NewTransactionRepositoryPostgres(db),
		CategoryRepository:          NewCategoryRepositoryPostgres(db),
		PrescribedExpanseRepository: NewPrescribedExpanseRepositoryPostgres(db),
	}
}
