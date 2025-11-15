package router

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"litespend-api/internal/httpsrv/middleware"
	"litespend-api/internal/model"
	"litespend-api/internal/service"
	"net/http"
	"strconv"
)

type PrescribedExpanseRouter struct {
	service *service.Service
}

func NewPrescribedExpanseRouter(service *service.Service) *PrescribedExpanseRouter {
	return &PrescribedExpanseRouter{
		service: service,
	}
}

func (r *PrescribedExpanseRouter) CreatePrescribedExpanse(c *gin.Context) {
	logined, ok := middleware.GetUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req model.CreatePrescribedExpanseRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := r.service.PrescribedExpanse.Create(c.Request.Context(), logined, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (r *PrescribedExpanseRouter) UpdatePrescribedExpanse(c *gin.Context) {
	logined, ok := middleware.GetUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid prescribed expanse id"})
		return
	}

	var req model.UpdatePrescribedExpanseRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = r.service.PrescribedExpanse.Update(c.Request.Context(), logined, id, req)
	if err != nil {
		if errors.Is(err, service.ErrPrescribedExpanseNotFound) {
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

	c.JSON(http.StatusOK, gin.H{"message": "prescribed expanse updated"})
}

func (r *PrescribedExpanseRouter) DeletePrescribedExpanse(c *gin.Context) {
	logined, ok := middleware.GetUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid prescribed expanse id"})
		return
	}

	err = r.service.PrescribedExpanse.Delete(c.Request.Context(), logined, id)
	if err != nil {
		if errors.Is(err, service.ErrPrescribedExpanseNotFound) {
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

	c.JSON(http.StatusOK, gin.H{"message": "prescribed expanse deleted"})
}

func (r *PrescribedExpanseRouter) GetPrescribedExpanse(c *gin.Context) {
	logined, ok := middleware.GetUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid prescribed expanse id"})
		return
	}

	prescribedExpanse, err := r.service.PrescribedExpanse.GetByID(c.Request.Context(), logined, id)
	if err != nil {
		if errors.Is(err, service.ErrPrescribedExpanseNotFound) {
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

	c.JSON(http.StatusOK, prescribedExpanse)
}

func (r *PrescribedExpanseRouter) GetPrescribedExpanses(c *gin.Context) {
	logined, ok := middleware.GetUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	prescribedExpanses, err := r.service.PrescribedExpanse.GetList(c.Request.Context(), logined)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prescribedExpanses)
}

func (r *PrescribedExpanseRouter) GetPrescribedExpansesWithPaymentStatus(c *gin.Context) {
	logined, ok := middleware.GetUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	prescribedExpanses, err := r.service.PrescribedExpanse.GetListWithPaymentStatus(c.Request.Context(), logined)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prescribedExpanses)
}

func (r *PrescribedExpanseRouter) MarkAsPaid(c *gin.Context) {
	logined, ok := middleware.GetUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid prescribed expanse id"})
		return
	}

	transactionID, err := r.service.PrescribedExpanse.MarkAsPaid(c.Request.Context(), logined, id)
	if err != nil {
		if errors.Is(err, service.ErrPrescribedExpanseNotFound) {
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

	c.JSON(http.StatusOK, gin.H{"transaction_id": transactionID, "message": "prescribed expanse marked as paid"})
}

func (r *PrescribedExpanseRouter) MarkAsPaidPartial(c *gin.Context) {
	logined, ok := middleware.GetUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid prescribed expanse id"})
		return
	}

	var req struct {
		Amount string `json:"amount" binding:"required"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount format"})
		return
	}

	transactionID, err := r.service.PrescribedExpanse.MarkAsPaidPartial(c.Request.Context(), logined, id, amount)
	if err != nil {
		if errors.Is(err, service.ErrPrescribedExpanseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrAccessDenied) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction_id": transactionID, "message": "prescribed expanse partially marked as paid"})
}
