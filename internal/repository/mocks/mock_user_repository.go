package mocks

import (
	"github.com/samvibes/vexop/auth-service/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)

	user := &models.User{
		Email: email,
	}

	return user, args.Error(0)
}
