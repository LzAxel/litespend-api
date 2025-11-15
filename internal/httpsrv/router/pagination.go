package router

import (
	"github.com/gin-gonic/gin"
	"litespend-api/internal/model"
	"strconv"
)

func ParsePaginationFromContext(c *gin.Context) model.PaginationParams {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	return model.NewPaginationParams(page, limit)
}
