package router

import "litespend-api/internal/service"

type Router struct {
	User              *UserRouter
	Transaction       *TransactionRouter
	Category          *CategoryRouter
	PrescribedExpanse *PrescribedExpanseRouter
	Auth              *AuthRouter
}

func NewRouter(service *service.Service) *Router {
	return &Router{
		User:              NewUserRouter(service),
		Transaction:       NewTransactionRouter(service),
		Category:          NewCategoryRouter(service),
		PrescribedExpanse: NewPrescribedExpanseRouter(service),
		Auth:              NewAuthRouter(service),
	}
}
