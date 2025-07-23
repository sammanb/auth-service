package repository

import (
	"github.com/google/uuid"
	"github.com/samvibes/vexop/auth-service/internal/models"
	"github.com/samvibes/vexop/auth-service/internal/utils"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	CopyPermissionsTx(tx *gorm.DB, tenant_id string) (utils.PermissionMap, error)
}

type PermissionRepo struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &PermissionRepo{db: db}
}

func (p *PermissionRepo) CopyPermissionsTx(tx *gorm.DB, tenant_id string) (utils.PermissionMap, error) {
	var permissions []*models.Permission
	permissionMap := make(utils.PermissionMap)
	if err := p.db.Where("tenant_id IS NULL").Find(&permissions).Error; err != nil {
		return nil, err
	}
	tenantUUID, err := uuid.Parse(tenant_id)
	if err != nil {
		return nil, err
	}
	var newPermissions []*models.Permission
	for _, permission := range permissions {
		id := uuid.New()
		newPermission := &models.Permission{
			ID:       id,
			TenantID: &tenantUUID,
			Action:   permission.Action,
			Resource: permission.Resource,
			Code:     permission.Code,
		}
		permissionMap[permission.ID] = newPermission
		newPermissions = append(newPermissions, newPermission)
	}

	if err := p.db.Save(newPermissions).Error; err != nil {
		return nil, err
	}

	return permissionMap, nil
}
