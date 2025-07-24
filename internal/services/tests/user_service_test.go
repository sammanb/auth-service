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
	"github.com/stretchr/testify/require"
)

func TestFindUserByEmail_Success(t *testing.T) {
	mockUserRepo := &mocks.MockUserRepository{}
	mockRoleRepo := &mocks.MockRoleRepository{}
	mockPermissionRepo := &mocks.MockPermissionRepository{}
	authService := &mocks.MockAuthService{}
	userService := services.NewUserService(mockUserRepo, mockRoleRepo, mockPermissionRepo, authService)

	email := "testuser@mail.com"

	expectedUser := &models.User{
		Email: email,
	}

	mockUserRepo.On("FindUserByEmail", email).Return(expectedUser, nil)

	user, err := userService.FindUserByEmail(email)

	assert.NoError(t, err)
	mockUserRepo.AssertCalled(t, "FindUserByEmail", email)
	require.NotNil(t, user)
	assert.Equal(t, expectedUser.Email, user.Email)
}

func TestFindUserByEmail_Failure(t *testing.T) {
	mockUserRepo := &mocks.MockUserRepository{}
	mockRoleRepo := &mocks.MockRoleRepository{}
	mockPermissionRepo := &mocks.MockPermissionRepository{}
	authService := &mocks.MockAuthService{}
	userService := services.NewUserService(mockUserRepo, mockRoleRepo, mockPermissionRepo, authService)

	email := "testuser@mail.com"

	mockUserRepo.On("FindUserByEmail", email).Return(nil, errors.New("user not found"))

	user, err := userService.FindUserByEmail(email)

	assert.Error(t, err)
	mockUserRepo.AssertCalled(t, "FindUserByEmail", email)
	require.Nil(t, user)
}

func TestCreateUser_Success(t *testing.T) {
	db := utils.SetupTestDB(t)

	mockUserRepo := &mocks.MockUserRepository{}
	mockRoleRepo := &mocks.MockRoleRepository{}
	mockPermissionRepo := &mocks.MockPermissionRepository{}
	authService := &mocks.MockAuthService{}
	userService := services.NewUserService(mockUserRepo, mockRoleRepo, mockPermissionRepo, authService)

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
	authService := &mocks.MockAuthService{}
	userService := services.NewUserService(mockUserRepo, mockRoleRepo, mockPermissionRepo, authService)

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

func TestLogin_Success(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	permissionRepo := &mocks.MockPermissionRepository{}
	roleRepo := &mocks.MockRoleRepository{}
	authService := &mocks.MockAuthService{}
	userService := services.NewUserService(userRepo, roleRepo, permissionRepo, authService)

	email := "testuser@mail.com"
	password := "password"
	expectedUser := &models.User{
		Email:        email,
		PasswordHash: "PasswordHash",
	}

	token := "123"

	userRepo.On("FindUserByEmail", email).Return(expectedUser, nil)
	authService.On("CompareHashAndPassword", []byte(password), []byte(expectedUser.PasswordHash)).Return(true)
	authService.On("GenerateJWT", mock.Anything).Return(token, nil)

	result, err := userService.Login(email, password)

	assert.NoError(t, err)
	require.NotEmpty(t, result)
	authService.AssertCalled(t, "GenerateJWT", expectedUser)
	authService.AssertExpectations(t)
}
