package router

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"litespend-api/internal/httpsrv/middleware"
	"litespend-api/internal/model"
	"litespend-api/internal/service"
	"net/http"
)

type ImportRouter struct {
	service *service.Service
}

func NewImportRouter(service *service.Service) *ImportRouter {
	return &ImportRouter{
		service: service,
	}
}

func (r *ImportRouter) ParseExcelFile(c *gin.Context) {
	_, ok := middleware.GetUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to open file"})
		return
	}
	defer src.Close()

	fileData, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read file"})
		return
	}

	structure, err := r.service.Import.ParseExcelFile(fileData)
	if err != nil {
		if errors.Is(err, service.ErrInvalidFileFormat) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file format"})
			return
		}
		if errors.Is(err, service.ErrEmptyFile) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file is empty"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, structure)
}

func (r *ImportRouter) ImportData(c *gin.Context) {
	logined, ok := middleware.GetUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	var mapping model.ExcelColumnMapping
	// Парсим маппинг из JSON строки в form data
	mappingJSON := c.PostForm("mapping")
	if mappingJSON != "" {
		if err := json.Unmarshal([]byte(mappingJSON), &mapping); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mapping format"})
			return
		}
	} else {
		// Пробуем получить из JSON body
		if err := c.ShouldBindJSON(&mapping); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "mapping is required"})
			return
		}
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to open file"})
		return
	}
	defer src.Close()

	fileData, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read file"})
		return
	}

	result, err := r.service.Import.ImportData(c.Request.Context(), logined, fileData, mapping)
	if err != nil {
		if errors.Is(err, service.ErrInvalidFileFormat) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file format"})
			return
		}
		if errors.Is(err, service.ErrEmptyFile) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file is empty"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
