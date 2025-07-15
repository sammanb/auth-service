package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/samvibes/vexop/auth-service/internal/middleware"
	"github.com/samvibes/vexop/auth-service/internal/models"
)

func GetCurrentUser(c *gin.Context) *models.User {
	userVar, exists := c.Get(middleware.UserContextKey)
	if !exists {
		return nil
	}
	user := userVar.(models.User)
	return &user
}
