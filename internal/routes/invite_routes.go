package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samvibes/vexop/auth-service/internal/handlers"
)

func RegisterInviteRoutes(group *gin.RouterGroup, inviteHandler handlers.InviteHandler) {
	group.GET("/", inviteHandler.GetInvites)
	group.POST("/", inviteHandler.CreateInvite)
	group.DELETE("/", inviteHandler.RemoveInvite)
	group.PUT("/accept", inviteHandler.AcceptInvite)
	group.POST("/resend", inviteHandler.ResendInvitation)
}
