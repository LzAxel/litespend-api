package router

import (
	"litespend-api/internal/service"
	"litespend-api/internal/session"
)

type Router struct {
	User *UserRouter
}

func NewRouter(service *service.Service, sm *session.SessionManager) *Router {
	return &Router{
		User: NewUserRouter(service, sm),
	}
}
