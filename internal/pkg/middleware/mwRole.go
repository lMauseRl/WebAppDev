package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lud0m4n/WebAppDev/internal/http/repository"
	"github.com/lud0m4n/WebAppDev/internal/model"
)

func ModeratorOnly(r *repository.Repository, c *gin.Context) bool {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
		c.Abort()
	}

	userID := ctxUserID.(uint)

	role, err := r.GetUserRoleByID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
	}	

	if role == model.ModeratorRole {
		return true
	}
	return false
}
