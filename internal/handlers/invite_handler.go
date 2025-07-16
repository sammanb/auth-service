package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samvibes/vexop/auth-service/internal/dto"
	"github.com/samvibes/vexop/auth-service/internal/services"
	"github.com/samvibes/vexop/auth-service/internal/utils"
)

type InviteHandler struct {
	inviteService services.InviteService
}

func NewInviteHandler(inviteService services.InviteService) *InviteHandler {
	return &InviteHandler{inviteService: inviteService}
}

func (i *InviteHandler) CreateInvite(c *gin.Context) {
	var createInviteReq dto.CreateInviteRequest

	if err := c.ShouldBindJSON(&createInviteReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	requestor := utils.GetCurrentUser(c)
	err := i.inviteService.CreateInvite(requestor, createInviteReq.Email, createInviteReq.Role, createInviteReq.TenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "invite sent"})
}

func (i *InviteHandler) GetInvites(c *gin.Context) {

}

func (i *InviteHandler) RemoveInvite(c *gin.Context) {}

func (i *InviteHandler) AcceptInvite(c *gin.Context) {}
