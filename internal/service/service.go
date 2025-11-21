package service

import (
	"context"
	"litespend-api/internal/model"
	"litespend-api/internal/repository"
	"litespend-api/internal/session"
	"time"
)

type Service struct {
	User
	Transaction
	Category
	Budget
	Auth
	Import
}

type User interface {
	Register(ctx context.Context, user model.RegisterRequest) error
	Login(ctx context.Context, req model.LoginRequest) (model.User, error)
	Logout(ctx context.Context) error
}

type Transaction interface {
	Create(ctx context.Context, logined model.User, transaction model.CreateTransactionRequest) (int, error)
	Update(ctx context.Context, logined model.User, id int, dto model.UpdateTransactionRequest) error
	Delete(ctx context.Context, logined model.User, id int) error
	GetByID(ctx context.Context, logined model.User, id int) (model.Transaction, error)
	GetList(ctx context.Context, logined model.User) ([]model.Transaction, error)
	GetListPaginated(ctx context.Context, logined model.User, params model.PaginationParams) (model.PaginatedTransactionsResponse, error)
	GetBalanceStatistics(ctx context.Context, logined model.User) (model.CurrentBalanceStatistics, error)
	GetCategoryStatistics(ctx context.Context, logined model.User, period model.PeriodType, from, to *time.Time) (model.CategoryStatisticsResponse, error)
	GetPeriodStatistics(ctx context.Context, logined model.User, period model.PeriodType, from, to *time.Time) (model.PeriodStatisticsResponse, error)
}

type Category interface {
	Create(ctx context.Context, logined model.User, req model.CreateCategoryRequest) (int, error)
	Update(ctx context.Context, logined model.User, id int, dto model.UpdateCategoryRequest) error
	Delete(ctx context.Context, logined model.User, id int) error
	GetByID(ctx context.Context, logined model.User, id int) (model.TransactionCategory, error)
	GetList(ctx context.Context, logined model.User) ([]model.TransactionCategory, error)
	GetListByType(ctx context.Context, logined model.User, categoryType model.CategoryType) ([]model.TransactionCategory, error)
}

type Budget interface {
	Create(ctx context.Context, logined model.User, req model.CreateBudgetRequest) (int, error)
	Update(ctx context.Context, logined model.User, id int, dto model.UpdateBudgetRequest) error
	Delete(ctx context.Context, logined model.User, id int) error
	GetByID(ctx context.Context, logined model.User, id int) (model.Budget, error)
	GetList(ctx context.Context, logined model.User) ([]model.Budget, error)
	GetListByPeriod(ctx context.Context, logined model.User, year uint, month uint) ([]model.Budget, error)
}

type Auth interface {
	RevokeSession(ctx context.Context, logined model.User, token string) error
	GetSessionInfo(ctx context.Context, logined model.User, token string) (model.SessionInfo, error)
}

type Import interface {
	ParseExcelFile(fileData []byte) (model.ExcelFileStructure, error)
	ImportData(ctx context.Context, logined model.User, fileData []byte, mapping model.ExcelColumnMapping) (model.ImportResult, error)
}

func NewService(repository *repository.Repository, sessionManager *session.SessionManager) *Service {
	return &Service{
		User:        NewUserService(repository.UserRepository),
		Transaction: NewTransactionService(repository.TransactionRepository),
		Category:    NewCategoryService(repository.CategoryRepository),
		Budget:      NewBudgetService(repository.BudgetRepository),
		Auth:        NewAuthService(sessionManager, repository.UserRepository),
		Import:      NewImportService(repository.TransactionRepository, repository.CategoryRepository),
	}
}
