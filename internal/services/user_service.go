package services

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/samvibes/vexop/auth-service/internal/models"
	"github.com/samvibes/vexop/auth-service/internal/repository"
	"github.com/samvibes/vexop/auth-service/internal/utils"
	"gorm.io/gorm"
)

type UserService interface {
	FindUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User, db *gorm.DB) error
	Login(email, password string) (string, error)
	RemoveUserById(tenant_id, user_id string) error
	RemoveUserByEmail(tenant_id string, email string) error
	InitResetPassword(email string) (string, error)
	ResetPassword(tenant_id, user_id, token, password string) error
	GetUsers(tenant_id string, page, limit int) ([]*models.User, error)
	GetUserById(tenant_id, user_id string) (*models.User, error)
	UpdateUserRole(tenant_id, user_id, role_name string) error
}

type UserServiceImpl struct {
	userRepo        repository.UserRepository
	roleRepo        repository.RoleRepository
	permissionsRepo repository.PermissionRepository
	authService     AuthService
}

func NewUserService(
	repo repository.UserRepository,
	role repository.RoleRepository,
	permission repository.PermissionRepository,
	authService AuthService,
) UserService {
	return &UserServiceImpl{userRepo: repo, roleRepo: role, permissionsRepo: permission, authService: authService}
}

func (u *UserServiceImpl) FindUserByEmail(email string) (*models.User, error) {
	return u.userRepo.FindUserByEmail(email)
}

func (u *UserServiceImpl) CreateUser(user *models.User, db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		// copy roles and permissions with this user's tenant id
		permissionMap, err := u.permissionsRepo.CopyPermissionsTx(tx, user.TenantID.String())
		if err != nil {
			return err
		}
		roles, err := u.roleRepo.CopyRolesTx(tx, user.TenantID.String(), &permissionMap)
		if err != nil {
			return err
		}

		var defaultRole *models.Role
		for _, role := range roles {
			if role.IsDefault {
				defaultRole = role
				break
			}
		}

		user.Role = *defaultRole
		user.RoleID = defaultRole.ID.String()
		err = u.userRepo.CreateUserTx(tx, user)
		if utils.UniqueViolation(err) {
			return utils.NewAppError(http.StatusBadRequest, "user already exists")
		}
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (u *UserServiceImpl) Login(email, password string) (string, error) {
	// check if user exists
	user, err := u.userRepo.FindUserByEmail(email)
	if err != nil {
		return "", err
	}

	// match password
	if !u.authService.CompareHashAndPassword([]byte(password), []byte(user.PasswordHash)) {
		return "", errors.New("incorrect password")
	}

	// generate jwt token
	return u.authService.GenerateJWT(user)
}

func (u *UserServiceImpl) RemoveUserById(tenant_id, user_id string) error {
	if err := u.userRepo.RemoveUserById(tenant_id, user_id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appError := utils.NewAppError(http.StatusNotFound, "user not found")
			return appError
		}

		return err
	}
	return nil
}

func (u *UserServiceImpl) RemoveUserByEmail(tenant_id string, email string) error {
	if err := u.userRepo.RemoveUserByEmail(tenant_id, email); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appError := utils.NewAppError(http.StatusNotFound, "user not found")
			return appError
		}

		return err
	}
	return nil
}

func (u *UserServiceImpl) InitResetPassword(email string) (string, error) {
	user, err := u.userRepo.FindUserByEmail(email)
	if err != nil {
		return "", err
	}

	// generate token to send back
	token, hashToken, err := utils.CreateRandomToken()
	if err != nil {
		return "", fmt.Errorf("error while generating token %s", err.Error())
	}

	err = u.userRepo.SetResetPasswordTokenHash(user.ID.String(), hashToken)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserServiceImpl) ResetPassword(tenant_id, user_id, token, password string) error {
	user, err := u.userRepo.GetUserById(tenant_id, user_id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appErr := utils.NewAppError(http.StatusNotFound, "user not found")
			return appErr
		}

		return err
	}

	// check token against ResetPasswordTokenHash
	if !u.authService.CompareHashAndPassword([]byte(token), []byte(user.ResetPasswordTokenHash)) {
		appErr := utils.NewAppError(http.StatusBadRequest, "incorrect password reset token")
		return appErr
	}

	newPasswordHash, err := u.authService.HashPassword(password)
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

func (u *UserServiceImpl) GetUsers(tenant_id string, page, limit int) ([]*models.User, error) {
	return u.userRepo.GetUsers(tenant_id, page, limit)
}

func (u *UserServiceImpl) GetUserById(tenant_id, user_id string) (*models.User, error) {
	return u.userRepo.GetUserById(tenant_id, user_id)
}

func (u *UserServiceImpl) UpdateUserRole(tenant_id, user_id, role_name string) error {
	role, err := u.roleRepo.GetRoleByName(tenant_id, role_name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appError := utils.NewAppError(http.StatusNotFound, "role not found")
			return appError
		}
		return err
	}

	user, err := u.userRepo.GetUserById(tenant_id, user_id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appError := utils.NewAppError(http.StatusNotFound, "user not found")
			return appError
		}
		return err
	}

	user.Role = *role
	user.RoleID = role.ID.String()

	return u.userRepo.UpdateUser(user)
}
