package repository

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/samvibes/vexop/auth-service/internal/dto"
	"github.com/samvibes/vexop/auth-service/internal/models"
	"github.com/samvibes/vexop/auth-service/internal/utils"
	"gorm.io/gorm"
)

type RoleRepository interface {
	GetRoleByName(tenant_id, role_name string) (*models.Role, error)
	GetRoles(tenant_id string, page, limit int) ([]*dto.RoleResponse, error)
	AddRole(tenant_id, name string) error
	AddRolePermission(tenant_id, id string, permission *models.Permission) error
	RemoveRolePermission(tenant_id, id string, permission *models.Permission) error
	UpdateRolePermissions(role *models.Role) error
	CopyRolesTx(tx *gorm.DB, tenant_id string, permissionMap *utils.PermissionMap) ([]*models.Role, error)
}

type RoleRepo struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepo {
	return &RoleRepo{db: db}
}

func (r *RoleRepo) GetRoleByName(tenant_id, role_name string) (*models.Role, error) {
	var role models.Role
	if err := r.db.Where("tenant_id = ? AND name = ?", tenant_id, role_name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepo) GetRoles(tenant_id string, page, limit int) ([]*dto.RoleResponse, error) {
	offset := (page - 1) * limit

	var _roles []*models.Role

	if err := r.db.Where("tenant_id = ?", tenant_id).Find(&_roles).Offset(offset).Limit(limit).Error; err != nil {
		return nil, err
	}

	var roles []*dto.RoleResponse

	for _, role := range _roles {
		var permissions []dto.PermissionInfo
		for _, permission := range role.Permissions {
			_permission := dto.PermissionInfo{
				ID:       permission.ID.String(),
				Action:   permission.Action,
				Resource: permission.Resource,
				Code:     permission.Code,
			}
			permissions = append(permissions, _permission)
		}

		roleObj := dto.RoleResponse{
			ID:          role.ID.String(),
			Name:        role.Name,
			Permissions: permissions,
			IsDefault:   role.IsDefault,
		}

		roles = append(roles, &roleObj)
	}

	return roles, nil
}

func (r *RoleRepo) AddRole(tenant_id, name string) error {
	tenantId, err := uuid.Parse(tenant_id)
	if err != nil {
		appError := utils.NewAppError(http.StatusBadRequest, "invalid tenant id")
		return appError
	}

	if err := r.db.Create(&models.Role{Name: name, TenantID: &tenantId}).Error; err != nil {
		return err
	}

	return nil
}

func (r *RoleRepo) AddRolePermission(tenant_id, id string, permission *models.Permission) error {
	return nil
}

func (r *RoleRepo) RemoveRolePermission(tenant_id, id string, permission *models.Permission) error {
	return nil
}

func (r *RoleRepo) UpdateRolePermissions(role *models.Role) error {
	return nil
}

func (r *RoleRepo) CopyRolesTx(tx *gorm.DB, tenant_id string, permissionMap *utils.PermissionMap) ([]*models.Role, error) {
	var roles []*models.Role
	if err := r.db.Preload("Permissions").Where("tenant_id IS NULL AND Name != ?", utils.RoleSuperAdmin).Find(&roles).Error; err != nil {
		return nil, err
	}
	tenantUUID, err := uuid.Parse(tenant_id)
	if err != nil {
		return nil, err
	}
	var newRoles []*models.Role
	for _, role := range roles {
		oldPermissions := role.Permissions
		var newPermissions []*models.Permission
		for _, oldPerm := range oldPermissions {
			newPermissions = append(newPermissions, (*permissionMap)[oldPerm.ID])
		}
		id := uuid.New()
		newRole := &models.Role{
			ID:          id,
			TenantID:    &tenantUUID,
			Name:        role.Name,
			IsDefault:   role.IsDefault,
			Permissions: newPermissions,
		}
		newRoles = append(newRoles, newRole)
	}

	if err := tx.Create(&newRoles).Error; err != nil {
		return nil, err
	}

	return newRoles, nil
}
