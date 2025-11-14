package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samvibes/vexop/auth-service/internal/handlers"
)

func RegisterAPIRoutes(group *gin.RouterGroup, authHandler handlers.AuthHandler) {
	group.GET("/health", authHandler.Health)
	group.POST("/signup", authHandler.SignUp)
	group.POST("/login", authHandler.Login)
}
