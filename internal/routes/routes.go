package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samvibes/vexop/auth-service/app"
	"github.com/samvibes/vexop/auth-service/internal/middleware"
	"github.com/spf13/viper"
)

func InitRoutes(container *app.AppContainer) *gin.Engine {
	jwtSecret := []byte(viper.GetString("JWT_SECRET"))

	router := gin.Default()
	auth_api := router.Group("/api/auth")
	RegisterAPIRoutes(auth_api, container.AuthHandler)

	router.Use(middleware.JWTAuthMiddleware(container.DB, jwtSecret))
	router.Use(middleware.AutoRBAC())

	// Super admin APIs
	sa_api := router.Group("/api/sa")
	RegisterSARoutes(sa_api, container.TenantHandler)

	invite_api := router.Group("/api/invites")
	RegisterInviteRoutes(invite_api, container.InviteHandler)

	return router
}
