package router

import (
	"errors"
	"github.com/gin-gonic/gin"
	"litespend-api/internal/httpsrv/middleware"
	"litespend-api/internal/service"
	"net/http"
)

type AuthRouter struct {
	service *service.Service
}

func NewAuthRouter(service *service.Service) *AuthRouter {
	return &AuthRouter{
		service: service,
	}
}

func (r *AuthRouter) RevokeSession(c *gin.Context) {
	logined, ok := middleware.GetUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.service.Auth.RevokeSession(c.Request.Context(), logined, req.Token)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrSessionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "session revoked successfully"})
}

func (r *AuthRouter) GetSessionInfo(c *gin.Context) {
	logined, ok := middleware.GetUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token parameter is required"})
		return
	}

	sessionInfo, err := r.service.Auth.GetSessionInfo(c.Request.Context(), logined, token)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrSessionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sessionInfo)
}
