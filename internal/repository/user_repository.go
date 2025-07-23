package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/samvibes/vexop/auth-service/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	CreateUserTx(tx *gorm.DB, user *models.User) error
	FindUserByEmail(email string) (*models.User, error)
	FindUserByEmailAndTenant(email string, tenant_id string) (*models.User, error)
	RemoveUserById(tenant_id, user_id string) error
	RemoveUserByEmail(tenant_id string, email string) error
	SetResetPasswordTokenHash(id, tokenHash string) error
	GetUsers(tenant_id string, page, limit int) ([]*models.User, error)
	GetUserById(tenant_id, user_id string) (*models.User, error)
	UpdateUser(user *models.User) error
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
		return errors.New("tenant: " + err.Error())
	}
	user.ID = uuid.New()
	user.Role = role
	user.RoleID = role.ID.String()
	return u.db.Create(user).Error
}

func (u *UserRepo) CreateUserTx(tx *gorm.DB, user *models.User) error {
	return tx.Create(user).Error
}

func (u *UserRepo) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := u.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) FindUserByEmailAndTenant(email string, tenant_id string) (*models.User, error) {
	var user models.User
	if err := u.db.Where("email = ? AND tenant_id = ?", email, tenant_id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepo) RemoveUserById(tenant_id, user_id string) error {
	if err := u.db.Where("tenant_id = ? AND id = ?", tenant_id, user_id).Delete(&models.User{}).Error; err != nil {
		return err
	}

	return nil
}

func (u *UserRepo) RemoveUserByEmail(tenant_id string, email string) error {
	if err := u.db.Where("email = ? AND tenant_id = ?", email, tenant_id).Delete(&models.User{}).Error; err != nil {
		return err
	}

	return nil
}

func (u *UserRepo) SetResetPasswordTokenHash(id, tokenHash string) error {
	if err := u.db.Where("id = ?", id).First(&models.User{}).Update("reset_password_token_hash", tokenHash).Error; err != nil {
		return err
	}

	return nil
}

func (u *UserRepo) GetUsers(tenant_id string, page, limit int) ([]*models.User, error) {
	offset := (page - 1) * limit
	var users []*models.User
	if err := u.db.Preload("Role").Preload("Role.Permissions").Where("tenant_id = ?", tenant_id).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserRepo) GetUserById(tenant_id, user_id string) (*models.User, error) {
	var user *models.User

	if err := u.db.Where("tenant_id = ? AND id = ?", tenant_id, user_id).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepo) UpdateUser(user *models.User) error {
	return u.db.Save(user).Error
}
