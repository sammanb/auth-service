package services

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/samvibes/vexop/auth-service/internal/dto"
	"github.com/samvibes/vexop/auth-service/internal/models"
	"github.com/samvibes/vexop/auth-service/internal/repository"
	"github.com/samvibes/vexop/auth-service/internal/utils"
	"gorm.io/gorm"
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

func (u *UserService) RemoveUserById(tenant_id, user_id string) error {
	if err := u.userRepo.RemoveUserById(tenant_id, user_id); err != nil {
		if err == gorm.ErrRecordNotFound {
			appError := utils.NewAppError(http.StatusNotFound, "user not found")
			return appError
		}

		return err
	}
	return nil
}

func (u *UserService) RemoveUserByEmail(email string, tenant_id string) error {
	if err := u.userRepo.RemoveUserByEmail(email, tenant_id); err != nil {
		if err == gorm.ErrRecordNotFound {
			appError := utils.NewAppError(http.StatusNotFound, "user not found")
			return appError
		}

		return err
	}
	return nil
}

func (u *UserService) InitResetPassword(email string) (string, error) {
	user, err := u.userRepo.FindUserByEmail(email)
	if err != nil {
		return "", err
	}

	// generate token to send back
	token, hashToken, err := utils.GenerateRandomToken()
	if err != nil {
		return "", fmt.Errorf("error while generating token %s", err.Error())
	}

	err = u.userRepo.SetResetPasswordTokenHash(user.ID.String(), hashToken)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserService) ResetPassword(tenant_id, user_id, token, password string) error {
	user, err := u.userRepo.GetUserById(tenant_id, user_id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			appErr := utils.NewAppError(http.StatusNotFound, "user not found")
			return appErr
		}

		return err
	}

	// check token against ResetPasswordTokenHash
	if !CompareHashAndPassword([]byte(token), []byte(user.ResetPasswordTokenHash)) {
		appErr := utils.NewAppError(http.StatusBadRequest, "incorrect password reset token")
		return appErr
	}

	newPasswordHash, err := HashPassword(password)
	if err != nil {
		return err
	}

	user.PasswordHash = newPasswordHash
	user.ResetPasswordTokenHash = ""

	err = u.userRepo.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) GetUsers(tenant_id string, page, limit int) ([]*models.User, error) {
	return u.userRepo.GetUsers(tenant_id, page, limit)
}

func (u *UserService) GetUserById(tenant_id, user_id string) (*models.User, error) {
	return u.userRepo.GetUserById(tenant_id, user_id)
}

func (u *UserService) UpdateUserRole(tenant_id, user_id, role_name string) error {
	role, err := u.userRepo.GetRoleByName(tenant_id, role_name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			appError := utils.NewAppError(http.StatusNotFound, "role not found")
			return appError
		}
		return err
	}

	user, err := u.userRepo.GetUserById(tenant_id, user_id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			appError := utils.NewAppError(http.StatusNotFound, "user not found")
			return appError
		}
		return err
	}

	user.Role = *role
	user.RoleID = role.ID.String()

	return u.userRepo.UpdateUser(user)
}

func (u *UserService) GetRoles(tenant_id string, page, limit int) ([]*dto.RoleResponse, error) {
	return u.userRepo.GetRoles(tenant_id, page, limit)
}
