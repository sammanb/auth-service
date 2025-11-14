package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samvibes/vexop/auth-service/internal/dto"
	"github.com/samvibes/vexop/auth-service/internal/services"
	"github.com/samvibes/vexop/auth-service/internal/utils"
)

type TenantHandler interface {
	GetTenants(*gin.Context)
	GetTenantById(*gin.Context)
	CreateTenant(*gin.Context)
	DeleteTenant(*gin.Context)
}

type TenantHandlerImpl struct {
	service services.TenantService
}

func NewTenantHandler(service services.TenantService) TenantHandler {
	return &TenantHandlerImpl{service: service}
}

func (h *TenantHandlerImpl) GetTenants(c *gin.Context) {
	idStr, _ := c.GetQuery("id")

	if idStr != "" {
		h.GetTenantById(c)
		return
	}

	page, limit := utils.GetPageAndLimit(c)

	requestor := utils.GetCurrentUser(c)

	tenants, err := h.service.GetTenants(requestor, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get tenants"})
		return
	}

	c.JSON(http.StatusOK, tenants)
}

func (h *TenantHandlerImpl) GetTenantById(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required in query"})
		return
	}

	requestor := utils.GetCurrentUser(c)

	tenant, err := h.service.GetTenantById(requestor, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

func (h *TenantHandlerImpl) CreateTenant(c *gin.Context) {
	var req dto.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	requestor := utils.GetCurrentUser(c)

	tenant, err := h.service.CreateTenant(requestor, req.Name)
	if err != nil {
		if err == services.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create tenant"})
		return
	}

	c.JSON(http.StatusCreated, tenant)
}

func (h *TenantHandlerImpl) DeleteTenant(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required in query"})
		return
	}

	requestor := utils.GetCurrentUser(c)

	didDelete, err := h.service.DeleteTenantById(requestor, id)
	if !didDelete {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "tenant could not be deleted"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tenant deleted successfully"})
}
