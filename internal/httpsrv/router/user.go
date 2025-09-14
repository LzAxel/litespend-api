package router

import (
	"github.com/gin-gonic/gin"
	"litespend-api/internal/model"
	"litespend-api/internal/service"
	"litespend-api/internal/session"
	"net/http"
)

type UserRouter struct {
	service *service.Service
	sm      *session.SessionManager
}

func NewUserRouter(service *service.Service, sm *session.SessionManager) *UserRouter {
	return &UserRouter{
		service: service,
		sm:      sm,
	}
}

func (r UserRouter) Register(c *gin.Context) {
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

func (r UserRouter) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := r.service.User.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	err = r.sm.RenewToken(c)
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.JSON(http.StatusOK, "ok")
}

func (r UserRouter) GetProfile(c *gin.Context) {

}
