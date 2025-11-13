package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samvibes/vexop/auth-service/internal/services"
	"github.com/samvibes/vexop/auth-service/internal/utils"
)

type RoleHandler struct {
	roleService services.RoleService
}

func NewRoleHandler(roleService services.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

func (r *RoleHandler) GetRoles(c *gin.Context) {
	user := utils.GetCurrentUser(c)

	page, limit := utils.GetPageAndLimit(c)

	roles, err := r.roleService.GetRoles(user.TenantID.String(), page, limit)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching roles"})
		return
	}

	c.JSON(http.StatusOK, roles)
}

func (r *RoleHandler) AddRole(c *gin.Context) {
	user := utils.GetCurrentUser(c)

	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty role name"})
		return
	}

	err := r.roleService.AddRole(user.TenantID.String(), name)
	if err != nil {
		if appError, ok := err.(*utils.AppError); ok {
			c.JSON(appError.Code, appError.Message)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add role"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "added role successfully"})
}

func (r *RoleHandler) DeleteRole(c *gin.Context) {
	user := utils.GetCurrentUser(c)

	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty role id"})
		return
	}

	if err := r.roleService.DeleteRole(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON("message": "role deleted successfully")
}
