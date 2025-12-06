package router

import (
	"litespend-api/internal/service"
	"litespend-api/internal/session"
)

type Router struct {
	User        *UserRouter
	Transaction *TransactionRouter
	Category    *CategoryRouter
	Budget      *BudgetRouter
	Auth        *AuthRouter
	Account     *AccountRouter
}

func NewRouter(service *service.Service, sessionManager *session.SessionManager) *Router {
	return &Router{
		User:        NewUserRouter(service, sessionManager),
		Transaction: NewTransactionRouter(service),
		Category:    NewCategoryRouter(service),
		Budget:      NewBudgetRouter(service),
		Auth:        NewAuthRouter(service),
		Account:     NewAccountRouter(service),
	}
}
