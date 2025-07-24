package mocks

import (
	"github.com/samvibes/vexop/auth-service/internal/models"
	"github.com/stretchr/testify/mock"
)

type AuthService interface {
	GenerateJWT(user *models.User) (string, error)
}

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) CompareHashAndPassword(password, hashed []byte) bool {
	args := m.Called(password, hashed)
	return args.Bool(0)
}

func (m *MockAuthService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) GenerateJWT(user *models.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}
