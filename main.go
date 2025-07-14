package main

import (
	"github.com/gin-gonic/gin"
	"github.com/samvibes/vexop/auth-service/app"
	"github.com/samvibes/vexop/auth-service/internal/middleware"
	"github.com/samvibes/vexop/auth-service/internal/routes"
)

func main() {
	container := app.InitApp()

	router := gin.Default()

	api := router.Group("/api")
	routes.RegisterAPIRoutes(api, container.AuthHandler)

	sa_api := router.Group("/api/sa")
	sa_api.Use(middleware.JWTAuthMiddleware())
	routes.RegisterSARoutes(sa_api, container.TenantHandler)

	router.Run(":9000")
}
