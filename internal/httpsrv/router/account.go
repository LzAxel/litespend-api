package router

import (
	"litespend-api/internal/httpsrv/middleware"
	"litespend-api/internal/model"
	"litespend-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	logined, ok := middleware.GetUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	
	accounts, err := r.service.Account.GetList(c, logined)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accounts": accounts})
}

func (r AccountRouter) UpdateAccount(c *gin.Context) {
	logined, ok := middleware.GetUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	var req model.UpdateAccountRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = r.service.Account.Update(c.Request.Context(), logined, id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "account updated"})
}

func (r AccountRouter) DeleteAccount(c *gin.Context) {
	logined, ok := middleware.GetUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 24)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	err = r.service.Account.Delete(c.Request.Context(), logined, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "account deleted"})
}

func (r AccountRouter) CreateAccount(c *gin.Context) {
	logined, ok := middleware.GetUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req model.CreateAccountRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := r.service.Account.Create(c.Request.Context(), logined, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}
