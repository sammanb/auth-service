package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samvibes/vexop/auth-service/app"
	"github.com/samvibes/vexop/auth-service/internal/middleware"
	"github.com/spf13/viper"
)

func InitRoutes(container *app.AppContainer) *gin.Engine {
	router := gin.Default()

	auth_api := router.Group("/api/auth")
	RegisterAPIRoutes(auth_api, container.AuthHandler)

	jwtSecret := []byte(viper.GetString("JWT_SECRET"))

	sa_api := router.Group("/api/sa")
	sa_api.Use(middleware.JWTAuthMiddleware(container.DB, jwtSecret))
	RegisterSARoutes(sa_api, container.TenantHandler)

	return router
}
