package router

import (
	"github.com/gin-gonic/gin"
	"litespend-api/internal/model"
	"litespend-api/internal/service"
	"net/http"
)

type UserRouter struct {
	service *service.Service
}

func NewUserRouter(service *service.Service) *UserRouter {
	return &UserRouter{
		service: service,
	}
}

func (r *UserRouter) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := r.service.User.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, "ok")
}

func (r *UserRouter) Test(c *gin.Context) {
	c.JSON(http.StatusOK)
}
