package tests

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/samvibes/vexop/auth-service/internal/models"
	"github.com/samvibes/vexop/auth-service/internal/services"
	"github.com/samvibes/vexop/auth-service/internal/services/mocks"
	"github.com/samvibes/vexop/auth-service/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser_Success(t *testing.T) {
	db := utils.SetupTestDB(t)

	mockUserRepo := &mocks.MockUserRepository{}
	mockRoleRepo := &mocks.MockRoleRepository{}
	mockPermissionRepo := &mocks.MockPermissionRepository{}
	userService := services.NewUserService(mockUserRepo, mockRoleRepo, mockPermissionRepo)

	userID := uuid.New()
	tenantID := uuid.New()
	user := &models.User{
		ID:       userID,
		Email:    "test@email.com",
		TenantID: &tenantID,
	}

	permissionId := uuid.New()
	permissionMap := make(utils.PermissionMap)
	permissionMap[permissionId] = &models.Permission{
		ID:       permissionId,
		TenantID: &tenantID,
		Action:   "action",
		Resource: "resource",
		Code:     "resource:action",
	}

	roles := []*models.Role{
		{
			TenantID:  &tenantID,
			Name:      "role1",
			IsDefault: true,
		},
		{
			TenantID:  &tenantID,
			Name:      "role2",
			IsDefault: false,
		},
	}

	mockPermissionRepo.On("CopyPermissionsTx", mock.Anything, tenantID.String()).Return(permissionMap, nil)
	mockRoleRepo.On("CopyRolesTx", mock.Anything, tenantID.String(), &permissionMap).Return(roles, nil)
	mockUserRepo.On("CreateUserTx", mock.Anything, user).Return(nil)

	err := userService.CreateUser(user, db)

	assert.NoError(t, err)
	mockPermissionRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockUserRepo.AssertCalled(t, "CreateUserTx", mock.Anything, user)
}

func TestCreateUser_Failure(t *testing.T) {
	db := utils.SetupTestDB(t)

	mockUserRepo := &mocks.MockUserRepository{}
	mockRoleRepo := &mocks.MockRoleRepository{}
	mockPermissionRepo := &mocks.MockPermissionRepository{}
	userService := services.NewUserService(mockUserRepo, mockRoleRepo, mockPermissionRepo)

	userID := uuid.New()
	tenantID := uuid.New()
	user := &models.User{
		ID:       userID,
		Email:    "test@email.com",
		TenantID: &tenantID,
	}

	permissionId := uuid.New()
	permissionMap := make(utils.PermissionMap)
	permissionMap[permissionId] = &models.Permission{
		ID:       permissionId,
		TenantID: &tenantID,
		Action:   "action",
		Resource: "resource",
		Code:     "resource:action",
	}

	roles := []*models.Role{
		{
			TenantID:  &tenantID,
			Name:      "role1",
			IsDefault: true,
		},
		{
			TenantID:  &tenantID,
			Name:      "role2",
			IsDefault: false,
		},
	}

	mockPermissionRepo.On("CopyPermissionsTx", mock.Anything, tenantID.String()).Return(permissionMap, nil)
	mockRoleRepo.On("CopyRolesTx", mock.Anything, tenantID.String(), &permissionMap).Return(roles, nil)
	mockUserRepo.On("CreateUserTx", mock.Anything, user).Return(errors.New("user not found"))

	err := userService.CreateUser(user, db)

	assert.Error(t, err)

	mockUserRepo.AssertExpectations(t)
}
