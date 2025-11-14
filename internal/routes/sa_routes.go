package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samvibes/vexop/auth-service/internal/handlers"
)

func RegisterSARoutes(group *gin.RouterGroup, tenantHandler handlers.TenantHandler) {
	group.GET("/tenants", tenantHandler.GetTenants)
	group.POST("/tenants", tenantHandler.CreateTenant)
	group.DELETE("/tenants", tenantHandler.DeleteTenant)
}
