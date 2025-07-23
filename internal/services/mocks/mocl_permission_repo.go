package mocks

import (
	"github.com/samvibes/vexop/auth-service/internal/utils"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockPermissionRepository struct {
	mock.Mock
}

func (m *MockPermissionRepository) CopyPermissionsTx(tx *gorm.DB, tenant_id string) (utils.PermissionMap, error) {
	args := m.Called(tx, tenant_id)

	if permissions, ok := args.Get(0).(utils.PermissionMap); ok {
		return permissions, nil
	}

	return nil, args.Error(1)
}
