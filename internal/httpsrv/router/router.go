package router

import (
	"litespend-api/internal/service"
	"litespend-api/internal/session"
)

type Router struct {
	User              *UserRouter
	Transaction       *TransactionRouter
	Category          *CategoryRouter
	PrescribedExpanse *PrescribedExpanseRouter
	Auth              *AuthRouter
	Import            *ImportRouter
}

func NewRouter(service *service.Service, sessionManager *session.SessionManager) *Router {
	return &Router{
		User:              NewUserRouter(service, sessionManager),
		Transaction:       NewTransactionRouter(service),
		Category:          NewCategoryRouter(service),
		PrescribedExpanse: NewPrescribedExpanseRouter(service),
		Auth:              NewAuthRouter(service),
		Import:            NewImportRouter(service),
	}
}
