package router

import (
	"errors"
	"github.com/gin-gonic/gin"
	"litespend-api/internal/model"
	"litespend-api/internal/service"
	"litespend-api/internal/session"
	"net/http"
)

type UserRouter struct {
	service        *service.Service
	sessionManager *session.SessionManager
}

func NewUserRouter(service *service.Service, sessionManager *session.SessionManager) *UserRouter {
	return &UserRouter{
		service:        service,
		sessionManager: sessionManager,
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

	user, err := r.service.User.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if r.sessionManager != nil {
		err = r.sessionManager.RenewToken(c.Request, c.Writer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
			return
		}

		r.sessionManager.Put(c.Request, c.Writer, "user_id", int(user.ID))
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
	err := r.service.User.Logout(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if r.sessionManager != nil {
		err = r.sessionManager.Destroy(c.Request, c.Writer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to destroy session"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
