package services

import (
	"errors"

	"github.com/samvibes/vexop/auth-service/internal/models"
	"github.com/samvibes/vexop/auth-service/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{userRepo: repo}
}

func (u *UserService) FindUserByEmail(email string) (*models.User, error) {
	return u.userRepo.FindUserByEmail(email)
}

func (u *UserService) CreateUser(user *models.User) error {
	return u.userRepo.CreateUser(user)
}

func (u *UserService) Login(email, password string) (string, error) {
	// check if user exists
	user, err := u.userRepo.FindUserByEmail(email)
	if err != nil {
		return "", err
	}

	// match password
	if CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) {
		return "", errors.New("incorrect password")
	}

	// generate jwt token
	return GenerateJWT(user)
}

func (u *UserService) RemoveUserById(id string) error {
	if err := u.userRepo.RemoveUserById(id); err != nil {
		return err
	}
	return nil
}

func (u *UserService) RemoveUserByEmail(email string, tenant_id string) error {
	if err := u.userRepo.RemoveUserByEmail(email, tenant_id); err != nil {
		return err
	}
	return nil
}
