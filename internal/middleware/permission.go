package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/samvibes/vexop/auth-service/config"
	"github.com/samvibes/vexop/auth-service/internal/models"
)

func RequirePermission(action, resource string) gin.HandlerFunc {
	code := fmt.Sprintf("%s:%s", action, resource)

	return func(c *gin.Context) {
		userVar, exists := c.Get(UserContextKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		user := userVar.(models.User)

		if err := config.DB.Preload("Role.Permissions").First(&user, "id = ?", user.ID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "failed to load user role"})
			return
		}

		for _, p := range user.Role.Permissions {
			if p.Code == code {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
	}
}

func HasPermission(c *gin.Context, action, resource string) bool {
	userVar, exists := c.Get(UserContextKey)
	if !exists {
		return false
	}

	user := userVar.(models.User)

	if strings.ToLower(user.Role.Name) == "superadmin" ||
		strings.ToLower(user.Role.Name) == "admin" {
		return true
	}

	// preload permissions if not already done
	if len(user.Role.Permissions) == 0 {
		config.DB.Preload("Role.Permissions").First(&user, "id = ?", user.ID)
	}

	code := fmt.Sprintf("%s:%s", resource, action)
	for _, p := range user.Role.Permissions {
		if p.Code == code {
			return true
		}
	}

	return false
}
