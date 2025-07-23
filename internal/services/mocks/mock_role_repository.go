package mocks

import (
	"github.com/samvibes/vexop/auth-service/internal/dto"
	"github.com/samvibes/vexop/auth-service/internal/models"
	"github.com/samvibes/vexop/auth-service/internal/utils"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) GetRoleByName(tenant_id, role_name string) (*models.Role, error) {
	args := m.Called(tenant_id, role_name)

	if role, ok := args.Get(0).(*models.Role); ok {
		return role, nil
	}

	return nil, args.Error(1)
}

func (m *MockRoleRepository) GetRoles(tenant_id string, page, limit int) ([]*dto.RoleResponse, error) {
	args := m.Called(tenant_id, page, limit)

	if roles, ok := args.Get(0).([]*dto.RoleResponse); ok {
		return roles, nil
	}

	return nil, args.Error(1)
}

func (m *MockRoleRepository) AddRole(tenant_id, name string) error {
	args := m.Called(tenant_id, name)

	return args.Error(0)
}

func (m *MockRoleRepository) AddRolePermission(tenant_id, id string, permission *models.Permission) error {
	args := m.Called(tenant_id, id, permission)

	return args.Error(0)
}

func (m *MockRoleRepository) RemoveRolePermission(tenant_id, id string, permission *models.Permission) error {
	args := m.Called(tenant_id, id, permission)

	return args.Error(0)
}

func (m *MockRoleRepository) UpdateRolePermissions(role *models.Role) error {
	args := m.Called(role)

	return args.Error(0)
}

func (m *MockRoleRepository) CopyRolesTx(tx *gorm.DB, tenant_id string, permissionMap *utils.PermissionMap) ([]*models.Role, error) {
	args := m.Called(tx, tenant_id, permissionMap)

	if roles, ok := args.Get(0).([]*models.Role); ok {
		return roles, nil
	}

	return nil, args.Error(1)
}
