package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samvibes/vexop/auth-service/internal/handlers"
)

func RegisterRoleRoutes(router *gin.RouterGroup, roleHandler handlers.RoleHandler) {
	router.GET("/", roleHandler.GetRoles)
}
