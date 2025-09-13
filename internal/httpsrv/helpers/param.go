package helpers

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func ParamInt64(c *gin.Context, key string) int64 {
	value, _ := strconv.ParseInt(c.Param(key), 10, 64)
	return value
}
