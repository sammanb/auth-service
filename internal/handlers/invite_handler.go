package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samvibes/vexop/auth-service/internal/dto"
	"github.com/samvibes/vexop/auth-service/internal/services"
	"github.com/samvibes/vexop/auth-service/internal/utils"
	"gorm.io/gorm"
)

type InviteHandler interface {
	GetInvites(*gin.Context)
	CreateInvite(*gin.Context)
	RemoveInvite(*gin.Context)
	AcceptInvite(*gin.Context)
	ResendInvitation(*gin.Context)
}

type InviteHandlerImpl struct {
	inviteService services.InviteService
	db            *gorm.DB
}

func NewInviteHandler(inviteService services.InviteService, db *gorm.DB) InviteHandler {
	return &InviteHandlerImpl{inviteService: inviteService, db: db}
}

func (i *InviteHandlerImpl) CreateInvite(c *gin.Context) {
	var createInviteReq dto.CreateInviteRequest

	if err := c.ShouldBindJSON(&createInviteReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	requestor := utils.GetCurrentUser(c)
	token, inviteId, err := i.inviteService.CreateInvite(requestor, createInviteReq.Email, createInviteReq.Role)
	if err != nil {
		if appError, ok := err.(*utils.AppError); ok {
			c.JSON(appError.Code, gin.H{"error": appError.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// TODO: send token via email
	fmt.Println("To send email with token and id ", token, inviteId)

	c.JSON(http.StatusCreated, gin.H{"message": "invite sent"})
}

func (i *InviteHandlerImpl) GetInvites(c *gin.Context) {
	page, limit := utils.GetPageAndLimit(c)

	requestor := utils.GetCurrentUser(c)
	invitations, err := i.inviteService.GetInvites(requestor, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, invitations)
}

func (i *InviteHandlerImpl) RemoveInvite(c *gin.Context) {
	inviteId := c.Query("id")

	fmt.Println(inviteId)

	err := i.inviteService.RemoveInvite(inviteId)
	if err != nil {
		if appError, ok := err.(*utils.AppError); ok {
			c.JSON(appError.Code, appError.Message)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "invite removed successfully"})
}

func (i *InviteHandlerImpl) AcceptInvite(c *gin.Context) {
	var acceptInviteReq dto.AcceptInviteRequest

	if err := c.ShouldBindJSON(&acceptInviteReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := i.inviteService.AcceptInvite(acceptInviteReq, i.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (i *InviteHandlerImpl) ResendInvitation(c *gin.Context) {
	inviteId := c.Query("invite_id")

	invitation, err := i.inviteService.GetInviteById(inviteId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "could not fetch invitation")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("invitation sent out for %s", invitation.Email)})
}
