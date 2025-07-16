package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/samvibes/vexop/auth-service/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	FindUserByEmail(email string) (*models.User, error)
	RemoveUserById(id string) error
	RemoveUserByEmail(email string, tenant_id string) error
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepo{db: db}
}

func (u *UserRepo) CreateUser(user *models.User) error {
	// fetch default role - there MUST be one default role
	var role models.Role
	err := u.db.Where("is_default = true").First(&role).Error
	if err != nil {
		return errors.New("critical error: no default role found")
	}
	// Also check if tenant_id is valid
	var tenant models.Tenant
	err = u.db.Where("id = ?", user.TenantID).First(&tenant).Error
	if err != nil {
		return errors.New("tenant " + err.Error())
	}
	user.ID = uuid.New()
	user.Role = role
	user.RoleID = role.ID.String()
	return u.db.Create(user).Error
}

func (u *UserRepo) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := u.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) RemoveUserById(id string) error {
	var user models.User
	if err := u.db.Where("id = ?", id).First(&user).Error; err != nil {
		return err
	}

	if err := u.db.Delete(user).Error; err != nil {
		return err
	}

	return nil
}

func (u *UserRepo) RemoveUserByEmail(id string, tenant_id string) error {
	var user models.User
	if err := u.db.Where("id = ?, tenant_id = ?", id, tenant_id).First(&user).Error; err != nil {
		return err
	}

	if err := u.db.Delete(user).Error; err != nil {
		return err
	}

	return nil
}
