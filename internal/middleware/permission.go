package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/samvibes/vexop/auth-service/internal/models"
	"github.com/samvibes/vexop/auth-service/internal/utils"
	"gorm.io/gorm"
)

func GetCurrentUser(c *gin.Context) *models.User {
	userVar, exists := c.Get(utils.UserContextKey)
	if !exists {
		return nil
	}
	user := userVar.(models.User)
	return &user
}

func requirePermission(action, resource string, db *gorm.DB) gin.HandlerFunc {
	code := fmt.Sprintf("%s:%s", action, resource)

	return func(c *gin.Context) {
		user := GetCurrentUser(c)

		if err := db.Preload("Role.Permissions").First(&user, "id = ?", user.ID).Error; err != nil {
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

func hasPermission(c *gin.Context, action, resource string, db *gorm.DB) bool {
	user := GetCurrentUser(c)

	if strings.ToLower(user.Role.Name) == "superadmin" {
		return true
	}

	// preload permissions if not already done
	if len(user.Role.Permissions) == 0 {
		db.Preload("Role.Permissions").First(&user, "id = ?", user.ID)
	}

	code := fmt.Sprintf("%s:%s", resource, action)
	for _, p := range user.Role.Permissions {
		if p.Code == code {
			return true
		}
	}

	return false
}

func extractResource(c *gin.Context) string {
	path := c.FullPath()

	// split and return the first segment
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 2 {
		return utils.Singularize(parts[1])
	}
	return ""
}

func AutoRBAC(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestor := GetCurrentUser(c)
		if requestor == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
			return
		}
		action := utils.MethodToAction[c.Request.Method]
		resource := extractResource(c)

		if action == "" || resource == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "unable to determine action or resource"})
			return
		}

		if !hasPermission(c, action, resource, db) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		c.Next()
	}
}
