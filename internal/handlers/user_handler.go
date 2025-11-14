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

type UserHandler interface {
	GetUsers(*gin.Context)
	GetUserById(*gin.Context)
	UpdateUserRole(*gin.Context)
	SendResetPassword(*gin.Context)
	ResetPassword(*gin.Context)
	DeleteUser(c *gin.Context)
}

type UserHandlerImpl struct {
	userService services.UserService
	db          *gorm.DB
}

func NewUserHandler(userService services.UserService, db *gorm.DB) UserHandler {
	return &UserHandlerImpl{userService: userService, db: db}
}

func (u *UserHandlerImpl) GetUsers(c *gin.Context) {
	user := utils.GetCurrentUser(c)

	page, limit := utils.GetPageAndLimit(c)

	users, err := u.userService.GetUsers(user.TenantID.String(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (u *UserHandlerImpl) GetUserById(c *gin.Context) {
	user := utils.GetCurrentUser(c)

	user_id := c.Param("id")
	if user_id == "" {
		u.GetUsers(c)
		return
	}

	user, err := u.userService.GetUserById(user.TenantID.String(), user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (u *UserHandlerImpl) UpdateUserRole(c *gin.Context) {
	user := utils.GetCurrentUser(c)

	var req dto.UpdateUserRoleRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.ID.String() == req.UserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot update self role"})
		return
	}

	err = u.userService.UpdateUserRole(user.TenantID.String(), req.UserID, req.RoleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "role updated successfully"})
}

func (u *UserHandlerImpl) DeleteUser(c *gin.Context) {
	user := utils.GetCurrentUser(c)
	email := c.Param("email")
	id := c.Param("id")

	if id == "" && email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "either user id or email must be provided"})
		return
	}

	var err error

	if id != "" {
		err = u.userService.RemoveUserById(user.TenantID.String(), id)
	} else if email != "" {
		err = u.userService.RemoveUserByEmail(user.TenantID.String(), email)
	}

	if appErr, ok := err.(*utils.AppError); ok {
		c.JSON(appErr.Code, appErr.Message)
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

func (u *UserHandlerImpl) SendResetPassword(c *gin.Context) {
	user := utils.GetCurrentUser(c)
	token, err := u.userService.InitResetPassword(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed while initiating reset password"})
		return
	}

	// send token to user via mail
	fmt.Println(token)

	c.JSON(http.StatusOK, gin.H{"message": "reset password message sent"})
}

func (u *UserHandlerImpl) ResetPassword(c *gin.Context) {
	user := utils.GetCurrentUser(c)

	var req dto.ResetPasswordRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	err = u.userService.ResetPassword(user.TenantID.String(), user.ID.String(), req.Token, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
}
