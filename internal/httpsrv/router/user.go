package router

import (
	"errors"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.service.User.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user registered successfully"})
}

func (r *UserRouter) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := r.service.User.Login(c.Request.Context(), c, req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role.String(),
		},
	})
}

func (r *UserRouter) Logout(c *gin.Context) {
	err := r.service.User.Logout(c.Request.Context(), c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
