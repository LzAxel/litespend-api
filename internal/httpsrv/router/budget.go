package router

import (
	"errors"
	"github.com/gin-gonic/gin"
	"litespend-api/internal/httpsrv/middleware"
	"litespend-api/internal/model"
	"litespend-api/internal/service"
	"net/http"
	"strconv"
)

type BudgetRouter struct {
	service *service.Service
}

func NewBudgetRouter(service *service.Service) *BudgetRouter {
	return &BudgetRouter{service: service}
}

func (r *BudgetRouter) CreateBudget(c *gin.Context) {
	logined, ok := middleware.GetUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req model.CreateBudgetAllocationRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := r.service.Budget.Create(c.Request.Context(), logined, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (r *BudgetRouter) UpdateBudget(c *gin.Context) {
	logined, ok := middleware.GetUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid budget id"})
		return
	}

	var req model.UpdateBudgetAllocationRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = r.service.Budget.Update(c.Request.Context(), logined, id, req)
	if err != nil {
		if errors.Is(err, service.ErrBudgetNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrAccessDenied) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "budget updated"})
}

func (r *BudgetRouter) GetBudget(c *gin.Context) {
	logined, ok := middleware.GetUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid budget id"})
		return
	}

	budget, err := r.service.Budget.GetByID(c.Request.Context(), logined, id)
	if err != nil {
		if errors.Is(err, service.ErrBudgetNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrAccessDenied) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, budget)
}

func (r *BudgetRouter) GetBudgets(c *gin.Context) {
	logined, ok := middleware.GetUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	month, err := strconv.ParseUint(c.Query("month"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month"})
		return 
	}

	year, err := strconv.ParseUint(c.Query("year"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
		return 
	}

	budgets, err := r.service.Budget.GetList(c.Request.Context(), logined,year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, budgets)
}
