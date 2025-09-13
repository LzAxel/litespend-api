package router

import "litespend-api/internal/service"

type Router struct {
	User *UserRouter
}

func NewRouter(service *service.Service) *Router {
	return &Router{
		User: NewUserRouter(service),
	}
}
