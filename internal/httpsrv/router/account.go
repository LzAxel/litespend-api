package router

import (
	"github.com/gin-gonic/gin"
	"litespend-api/internal/service"
)

type AccountRouter struct {
	service *service.Service
}

func NewAccountRouter(service *service.Service) *AccountRouter {
	return &AccountRouter{
		service: service,
	}
}

func (r AccountRouter) GetAccounts(c *gin.Context) {

}

func (r AccountRouter) UpdateAccount(c *gin.Context) {

}

func (r AccountRouter) DeleteAccount(c *gin.Context) {

}

func (r AccountRouter) CreateAccount(c *gin.Context) {

}
