package mocks

import (
	"github.com/samvibes/vexop/auth-service/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (u *MockUserService) FindUserByEmail(email string) (*models.User, error) {
	args := u.Called(email)

	user := &models.User{}

	return user, args.Error(0)
}
