package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/samvibes/vexop/auth-service/internal/models"
)

var jwtSecret = []byte("secret-key")

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid authorization header"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}

		user := models.User{
			ID:   parseUUID(claims["id"]),
			Role: models.UserRole(claims["role"].(string)),
		}

		if claims["tenant_id"] != nil {
			tenantID := parseUUID(claims["tenant_id"])
			user.TenantID = &tenantID
		}

		c.Set("user", user)
		c.Next()
	}
}

func parseUUID(v interface{}) uuid.UUID {
	str, _ := v.(string)
	id, _ := uuid.Parse(str)
	return id
}
