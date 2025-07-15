package main

import (
	"github.com/gin-gonic/gin"
	"github.com/samvibes/vexop/auth-service/app"
	"github.com/samvibes/vexop/auth-service/internal/middleware"
	"github.com/samvibes/vexop/auth-service/internal/routes"
	"github.com/spf13/viper"
)

func main() {
	container := app.InitApp()

	router := gin.Default()

	api := router.Group("/api")
	routes.RegisterAPIRoutes(api, container.AuthHandler)

	jwtSecret := []byte(viper.GetString("JWT_SECRET"))

	sa_api := router.Group("/api/sa")
	sa_api.Use(middleware.JWTAuthMiddleware(container.DB, jwtSecret))
	routes.RegisterSARoutes(sa_api, container.TenantHandler)

	router.Run(":9000")
}
