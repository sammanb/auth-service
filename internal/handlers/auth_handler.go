package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samvibes/vexop/auth-service/internal/dto"
	"github.com/samvibes/vexop/auth-service/internal/models"
	"github.com/samvibes/vexop/auth-service/internal/services"
	"gorm.io/gorm"
)

type AuthHandler interface {
	Health(*gin.Context)
	SignUp(*gin.Context)
	Login(*gin.Context)
}

type AuthHandlerImpl struct {
	authService   services.AuthService
	userService   services.UserService
	tenantService services.TenantService
	db            *gorm.DB
}

func NewAuthHandler(authService services.AuthService, userService services.UserService, tenantService services.TenantService, db *gorm.DB) AuthHandler {
	return &AuthHandlerImpl{authService, userService, tenantService, db}
}

func (h *AuthHandlerImpl) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "auth service is up and running"})
}

func (h *AuthHandlerImpl) SignUp(c *gin.Context) {
	var req dto.SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := h.authService.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tenant, err := h.tenantService.CreateTenant(nil, req.Email)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{
		ID:           uuid.New(),
		TenantID:     tenant.ID,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		IsOwner:      true,
	}

	if err := h.userService.CreateUser(user, h.db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
}

func (h *AuthHandlerImpl) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
