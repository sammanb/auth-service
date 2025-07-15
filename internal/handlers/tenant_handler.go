package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/samvibes/vexop/auth-service/internal/services"
)

type TenantHandler struct {
	service services.TenantSvcInterface
}

func NewTenantHandler(service services.TenantSvcInterface) *TenantHandler {
	return &TenantHandler{service: service}
}

type CreateTenantRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *TenantHandler) GetTenants(c *gin.Context) {
	requestor := GetCurrentUser(c)
	if requestor == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	pageStr, _ := c.GetQuery("page")
	limitStr, _ := c.GetQuery("limit")
	idStr, _ := c.GetQuery("id")

	if idStr != "" {
		h.GetTenantById(c)
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	tenants, err := h.service.GetTenants(requestor, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get tenants"})
		return
	}

	c.JSON(http.StatusOK, tenants)
}

func (h *TenantHandler) GetTenantById(c *gin.Context) {
	requestor := GetCurrentUser(c)
	if requestor == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	id, ok := c.GetQuery("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required in query"})
		return
	}

	tenant, err := h.service.GetTenantById(requestor, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

func (h *TenantHandler) CreateTenant(c *gin.Context) {
	requestor := GetCurrentUser(c)
	if requestor == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	var req CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

func (h *TenantHandler) DeleteTenant(c *gin.Context) {
	requestor := GetCurrentUser(c)
	if requestor == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	id, ok := c.GetQuery("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required in query"})
		return
	}

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
