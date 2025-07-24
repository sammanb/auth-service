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
	mockAuthService := &mocks.MockAuthService{}
	userService := services.NewUserService(mockUserRepo, mockRoleRepo, mockPermissionRepo, mockAuthService)

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
	mockAuthService := &mocks.MockAuthService{}
	userService := services.NewUserService(mockUserRepo, mockRoleRepo, mockPermissionRepo, mockAuthService)

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
	mockAuthService := &mocks.MockAuthService{}
	userService := services.NewUserService(mockUserRepo, mockRoleRepo, mockPermissionRepo, mockAuthService)

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
	mockUserRepo := &mocks.MockUserRepository{}
	mockRoleRepo := &mocks.MockRoleRepository{}
	mockPermissionRepo := &mocks.MockPermissionRepository{}
	mockAuthService := &mocks.MockAuthService{}
	userService := services.NewUserService(mockUserRepo, mockRoleRepo, mockPermissionRepo, mockAuthService)

	email := "testuser@mail.com"
	password := "password"
	expectedUser := &models.User{
		Email:        email,
		PasswordHash: "PasswordHash",
	}

	token := "123"

	mockUserRepo.On("FindUserByEmail", email).Return(expectedUser, nil)
	mockAuthService.On("CompareHashAndPassword", []byte(password), []byte(expectedUser.PasswordHash)).Return(true)
	mockAuthService.On("GenerateJWT", mock.Anything).Return(token, nil)

	result, err := userService.Login(email, password)

	assert.NoError(t, err)
	require.NotEmpty(t, result)
	mockAuthService.AssertCalled(t, "GenerateJWT", expectedUser)
	mockAuthService.AssertExpectations(t)
}

func TestRemoveUserById_Success(t *testing.T) {
	mockUserRepo := &mocks.MockUserRepository{}
	mockPermissionRepo := &mocks.MockPermissionRepository{}
	mockRoleRepo := &mocks.MockRoleRepository{}
	mockAuthService := &mocks.MockAuthService{}
	userService := services.NewUserService(mockUserRepo, mockRoleRepo, mockPermissionRepo, mockAuthService)

	tenant_id := "tenant_id"
	user_id := "user_id"
	mockUserRepo.On("RemoveUserById", tenant_id, user_id).Return(nil)

	err := userService.RemoveUserById(tenant_id, user_id)

	assert.Nil(t, err)
	mockUserRepo.AssertCalled(t, "RemoveUserById", tenant_id, user_id)
	mockUserRepo.AssertExpectations(t)
}

func TestRemoveUserById_Failure(t *testing.T) {
	mockUserRepo := &mocks.MockUserRepository{}
	mockPermissionRepo := &mocks.MockPermissionRepository{}
	mockRoleRepo := &mocks.MockRoleRepository{}
	mockAuthService := &mocks.MockAuthService{}
	userService := services.NewUserService(mockUserRepo, mockRoleRepo, mockPermissionRepo, mockAuthService)

	tenant_id := "tenant_id"
	user_id := "user_id"
	mockUserRepo.On("RemoveUserById", tenant_id, user_id).Return(errors.New("user not found"))

	err := userService.RemoveUserById(tenant_id, user_id)

	assert.NotNil(t, err)
	mockUserRepo.AssertCalled(t, "RemoveUserById", tenant_id, user_id)
	mockUserRepo.AssertExpectations(t)
}

func TestInitResetPassword_Success(t *testing.T) {
	mockUserRepo := &mocks.MockUserRepository{}
	mockPermissionRepo := &mocks.MockPermissionRepository{}
	mockRoleRepo := &mocks.MockRoleRepository{}
	mockAuthService := &mocks.MockAuthService{}
	userService := services.NewUserService(mockUserRepo, mockRoleRepo, mockPermissionRepo, mockAuthService)

	email := "testuser@mail.com"
	userId := uuid.New()

	expectedUser := &models.User{
		ID:    userId,
		Email: email,
	}

	mockUserRepo.On("FindUserByEmail", email).Return(expectedUser, nil)

	rawToken := "rawToken"
	hashedToken := "hashedToken"
	utils.CreateRandomToken = func() (string, string, error) {
		return rawToken, hashedToken, nil
	}
	defer func() {
		utils.CreateRandomToken = utils.GenerateRandomToken
	}()

	mockUserRepo.On("SetResetPasswordTokenHash", userId.String(), hashedToken).Return(nil)

	token, err := userService.InitResetPassword(email)

	assert.Nil(t, err)
	assert.Equal(t, rawToken, token)
	mockUserRepo.AssertCalled(t, "SetResetPasswordTokenHash", userId.String(), hashedToken)
	mockUserRepo.AssertExpectations(t)
}

func TestInitResetPassword_Failure(t *testing.T) {
	mockUserRepo := &mocks.MockUserRepository{}
	mockPermissionRepo := &mocks.MockPermissionRepository{}
	mockRoleRepo := &mocks.MockRoleRepository{}
	mockAuthService := &mocks.MockAuthService{}
	userService := services.NewUserService(mockUserRepo, mockRoleRepo, mockPermissionRepo, mockAuthService)

	email := "testuser@mail.com"
	userId := uuid.New()

	expectedUser := &models.User{
		ID:    userId,
		Email: email,
	}

	mockUserRepo.On("FindUserByEmail", email).Return(expectedUser, nil)

	rawToken := "rawToken"
	hashedToken := "hashedToken"
	utils.CreateRandomToken = func() (string, string, error) {
		return rawToken, hashedToken, nil
	}
	defer func() {
		utils.CreateRandomToken = utils.GenerateRandomToken
	}()

	mockUserRepo.On("SetResetPasswordTokenHash", userId.String(), hashedToken).Return(errors.New("error"))

	token, err := userService.InitResetPassword(email)

	assert.NotNil(t, err)
	assert.Equal(t, "", token)
	mockUserRepo.AssertCalled(t, "SetResetPasswordTokenHash", userId.String(), hashedToken)
	mockUserRepo.AssertExpectations(t)
}

func TestResetPassword_Success(t *testing.T) {
	mockUserRepo := &mocks.MockUserRepository{}
	mockPermissionRepo := &mocks.MockPermissionRepository{}
	mockRoleRepo := &mocks.MockRoleRepository{}
	mockAuthService := &mocks.MockAuthService{}
	userService := services.NewUserService(mockUserRepo, mockRoleRepo, mockPermissionRepo, mockAuthService)

	tenantId := uuid.New()
	userId := uuid.New()
	token := "token"
	newPassword := "newPassword"
	resetPasswordTokenHash := "resetPasswordTokenHash"

	user := &models.User{
		ID:                     userId,
		TenantID:               &tenantId,
		PasswordHash:           "hashedPassword",
		ResetPasswordTokenHash: resetPasswordTokenHash,
	}

	mockUserRepo.On("GetUserById", tenantId.String(), userId.String()).Return(user, nil)
	mockAuthService.On("CompareHashAndPassword", []byte(token), []byte(user.ResetPasswordTokenHash)).Return(true)
	mockAuthService.On("HashPassword", newPassword).Return("hashedPassword", nil)
	mockUserRepo.On("UpdateUser", user).Return(nil)

	err := userService.ResetPassword(tenantId.String(), userId.String(), token, newPassword)

	assert.Nil(t, err)
	mockUserRepo.AssertCalled(t, "GetUserById", tenantId.String(), userId.String())
	mockAuthService.AssertCalled(t, "CompareHashAndPassword", []byte(token), []byte(resetPasswordTokenHash))
	mockAuthService.AssertCalled(t, "HashPassword", newPassword)
	mockUserRepo.AssertCalled(t, "UpdateUser", user)

	mockUserRepo.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}

func TestResetPassword_InvalidUser_Failed(t *testing.T) {
	mockUserRepo := &mocks.MockUserRepository{}
	mockPermissionRepo := &mocks.MockPermissionRepository{}
	mockRoleRepo := &mocks.MockRoleRepository{}
	mockAuthService := &mocks.MockAuthService{}
	userService := services.NewUserService(mockUserRepo, mockRoleRepo, mockPermissionRepo, mockAuthService)

	tenantId := uuid.New()
	userId := uuid.New()
	token := "token"
	newPassword := "newPassword"

	mockUserRepo.On("GetUserById", tenantId.String(), userId.String()).Return(nil, errors.New("user not found"))

	err := userService.ResetPassword(tenantId.String(), userId.String(), token, newPassword)

	assert.NotNil(t, err)
	mockUserRepo.AssertCalled(t, "GetUserById", tenantId.String(), userId.String())
	mockUserRepo.AssertExpectations(t)
}

func TestResetPassword_InvalidToken_Failed(t *testing.T) {
	mockUserRepo := &mocks.MockUserRepository{}
	mockPermissionRepo := &mocks.MockPermissionRepository{}
	mockRoleRepo := &mocks.MockRoleRepository{}
	mockAuthService := &mocks.MockAuthService{}
	userService := services.NewUserService(mockUserRepo, mockRoleRepo, mockPermissionRepo, mockAuthService)

	tenantId := uuid.New()
	userId := uuid.New()
	token := "token"
	newPassword := "newPassword"
	resetPasswordTokenHash := "resetPasswordTokenHash"

	user := &models.User{
		ID:                     userId,
		TenantID:               &tenantId,
		PasswordHash:           "hashedPassword",
		ResetPasswordTokenHash: resetPasswordTokenHash,
	}

	mockUserRepo.On("GetUserById", tenantId.String(), userId.String()).Return(user, nil)
	mockAuthService.On("CompareHashAndPassword", []byte(token), []byte(user.ResetPasswordTokenHash)).Return(false)

	err := userService.ResetPassword(tenantId.String(), userId.String(), token, newPassword)

	assert.NotNil(t, err)
	mockUserRepo.AssertCalled(t, "GetUserById", tenantId.String(), userId.String())
	mockAuthService.AssertCalled(t, "CompareHashAndPassword", []byte(token), []byte(resetPasswordTokenHash))

	mockUserRepo.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}

func TestResetPassword_UpdateUser_Failure(t *testing.T) {
	mockUserRepo := &mocks.MockUserRepository{}
	mockPermissionRepo := &mocks.MockPermissionRepository{}
	mockRoleRepo := &mocks.MockRoleRepository{}
	mockAuthService := &mocks.MockAuthService{}
	userService := services.NewUserService(mockUserRepo, mockRoleRepo, mockPermissionRepo, mockAuthService)

	tenantId := uuid.New()
	userId := uuid.New()
	token := "token"
	newPassword := "newPassword"
	resetPasswordTokenHash := "resetPasswordTokenHash"

	user := &models.User{
		ID:                     userId,
		TenantID:               &tenantId,
		PasswordHash:           "hashedPassword",
		ResetPasswordTokenHash: resetPasswordTokenHash,
	}

	mockUserRepo.On("GetUserById", tenantId.String(), userId.String()).Return(user, nil)
	mockAuthService.On("CompareHashAndPassword", []byte(token), []byte(user.ResetPasswordTokenHash)).Return(true)
	mockAuthService.On("HashPassword", newPassword).Return("hashedPassword", nil)
	mockUserRepo.On("UpdateUser", user).Return(errors.New("user update failed"))

	err := userService.ResetPassword(tenantId.String(), userId.String(), token, newPassword)

	assert.NotNil(t, err)
	mockUserRepo.AssertCalled(t, "GetUserById", tenantId.String(), userId.String())
	mockAuthService.AssertCalled(t, "CompareHashAndPassword", []byte(token), []byte(resetPasswordTokenHash))
	mockAuthService.AssertCalled(t, "HashPassword", newPassword)
	mockUserRepo.AssertCalled(t, "UpdateUser", user)

	mockUserRepo.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}

func TestGetUsers_Success(t *testing.T) {
	mockUserRepo := &mocks.MockUserRepository{}
	mockPermissionRepo := &mocks.MockPermissionRepository{}
	mockRoleRepo := &mocks.MockRoleRepository{}
	mockAuthService := &mocks.MockAuthService{}
	userService := services.NewUserService(mockUserRepo, mockRoleRepo, mockPermissionRepo, mockAuthService)

	tenantId := uuid.New()
	userId := uuid.New()

	expectedUser := &models.User{
		ID:       userId,
		TenantID: &tenantId,
	}

	mockUserRepo.On("GetUserById", tenantId.String(), userId.String()).Return(expectedUser, nil)

	user, err := userService.GetUserById(tenantId.String(), userId.String())

	assert.Nil(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	mockUserRepo.AssertCalled(t, "GetUserById", tenantId.String(), userId.String())
}

func TestGetUsers_Failuer(t *testing.T) {
	mockUserRepo := &mocks.MockUserRepository{}
	mockPermissionRepo := &mocks.MockPermissionRepository{}
	mockRoleRepo := &mocks.MockRoleRepository{}
	mockAuthService := &mocks.MockAuthService{}
	userService := services.NewUserService(mockUserRepo, mockRoleRepo, mockPermissionRepo, mockAuthService)

	tenantId := uuid.New()
	userId := uuid.New()

	mockUserRepo.On("GetUserById", tenantId.String(), userId.String()).Return(nil, errors.New("user not found"))

	user, err := userService.GetUserById(tenantId.String(), userId.String())

	assert.NotNil(t, err)
	assert.Nil(t, user)
	mockUserRepo.AssertCalled(t, "GetUserById", tenantId.String(), userId.String())
}

func TestUpdateUserRole_Success(t *testing.T) {
	mockUserRepo := &mocks.MockUserRepository{}
	mockPermissionRepo := &mocks.MockPermissionRepository{}
	mockRoleRepo := &mocks.MockRoleRepository{}
	mockAuthService := &mocks.MockAuthService{}
	userService := services.NewUserService(mockUserRepo, mockRoleRepo, mockPermissionRepo, mockAuthService)

	tenantId := uuid.New()
	userId := uuid.New()
	roleName := "roleName"

	expectedRole := &models.Role{
		ID:       uuid.New(),
		TenantID: &tenantId,
	}

	user := &models.User{
		ID:       userId,
		TenantID: &tenantId,
		RoleID:   expectedRole.ID.String(),
		Role:     *expectedRole,
	}

	mockRoleRepo.On("GetRoleByName", tenantId.String(), roleName).Return(expectedRole, nil)
	mockUserRepo.On("GetUserById", tenantId.String(), userId.String()).Return(user, nil)
	mockUserRepo.On("UpdateUser", user).Return(nil)

	err := userService.UpdateUserRole(tenantId.String(), userId.String(), roleName)

	assert.Nil(t, err)
	mockRoleRepo.AssertCalled(t, "GetRoleByName", tenantId.String(), roleName)
	mockUserRepo.AssertCalled(t, "GetUserById", tenantId.String(), userId.String())
	mockUserRepo.AssertCalled(t, "UpdateUser", user)

	mockRoleRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}
