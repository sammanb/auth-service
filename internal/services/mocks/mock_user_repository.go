package mocks

import (
	"github.com/google/uuid"
	"github.com/samvibes/vexop/auth-service/internal/models"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) CreateUserTx(tx *gorm.DB, user *models.User) error {
	args := m.Called(tx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)

	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockUserRepository) FindUserByEmailAndTenant(email string, tenant_id string) (*models.User, error) {
	args := m.Called(email, tenant_id)

	tenantUUID, err := uuid.Parse(tenant_id)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    email,
		TenantID: &tenantUUID,
	}

	return user, args.Error(0)
}

func (m *MockUserRepository) RemoveUserById(tenant_id, user_id string) error {
	args := m.Called(tenant_id, user_id)

	return args.Error(0)
}

func (m *MockUserRepository) RemoveUserByEmail(tenant_id, email string) error {
	args := m.Called(tenant_id, email)

	return args.Error(0)
}

func (m *MockUserRepository) SetResetPasswordTokenHash(id, tokenHash string) error {
	args := m.Called(id, tokenHash)

	return args.Error(0)
}

func (m *MockUserRepository) GetUsers(tenant_id string, page int, limit int) ([]*models.User, error) {
	args := m.Called(tenant_id, page, limit)

	users := make([]*models.User, 0)
	users = append(users, &models.User{Email: "test_user@mail.com"})

	return users, args.Error(0)
}

func (m *MockUserRepository) GetUserById(tenant_id, user_id string) (*models.User, error) {
	args := m.Called(tenant_id, user_id)

	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockUserRepository) UpdateUser(user *models.User) error {
	args := m.Called(user)

	return args.Error(0)
}
