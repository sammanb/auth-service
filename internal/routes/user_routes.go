package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samvibes/vexop/auth-service/internal/handlers"
)

func RegisterUserRoutes(router *gin.RouterGroup, userHandler handlers.UserHandlerInterface) {
	router.POST("/password/forgot", userHandler.SendResetPassword)
	router.POST("/password/reset", userHandler.ResetPassword)
	router.GET("/", userHandler.GetUsers)
	router.GET("/:id", userHandler.GetUserById)
	router.PUT("/role", userHandler.UpdateUserRole)
	router.DELETE("/:id", userHandler.DeleteUser)
}
