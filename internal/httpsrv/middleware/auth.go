package middleware

import (
	"github.com/gin-gonic/gin"
	"litespend-api/internal/model"
	"litespend-api/internal/repository"
	"litespend-api/internal/session"
	"net/http"
	"strconv"
)

const UserContextKey = "user"

func RequireAuth(sessionManager *session.SessionManager, userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDVal := sessionManager.Get(c.Request, "user_id")
		if userIDVal == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		var userID uint64
		switch v := userIDVal.(type) {
		case uint64:
			userID = v
		case int:
			userID = uint64(v)
		case int64:
			userID = uint64(v)
		case string:
			parsed, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id in session"})
				c.Abort()
				return
			}
			userID = parsed
		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id type in session"})
			c.Abort()
			return
		}

		user, err := userRepo.GetByID(c.Request.Context(), int(userID))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			c.Abort()
			return
		}

		c.Set(UserContextKey, user)
		c.Next()
	}
}

func GetUserFromContext(c *gin.Context) (model.User, bool) {
	user, exists := c.Get(UserContextKey)
	if !exists {
		return model.User{}, false
	}

	userModel, ok := user.(model.User)
	if !ok {
		return model.User{}, false
	}

	return userModel, true
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := GetUserFromContext(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		if user.Role != model.UserRoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
