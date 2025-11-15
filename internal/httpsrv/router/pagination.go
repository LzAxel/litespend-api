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

	params := model.NewPaginationParams(page, limit)

	// Сортировка
	if sortByStr := c.Query("sort_by"); sortByStr != "" {
		sortBy := model.SortField(sortByStr)
		if sortBy == model.SortFieldDate || sortBy == model.SortFieldDescription || sortBy == model.SortFieldCategory {
			params.SortBy = &sortBy
		}
	}

	if sortOrderStr := c.Query("sort_order"); sortOrderStr != "" {
		sortOrder := model.SortOrder(sortOrderStr)
		if sortOrder == model.SortOrderASC || sortOrder == model.SortOrderDESC {
			params.SortOrder = &sortOrder
		}
	}

	// Поиск
	if searchStr := c.Query("search"); searchStr != "" {
		params.Search = &searchStr
	}

	return params
}
