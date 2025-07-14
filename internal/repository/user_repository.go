package repository

import (
	"github.com/google/uuid"
	"github.com/samvibes/vexop/auth-service/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	FindUserByEmail(email string) (*models.User, error)
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepo{db: db}
}

func (u *UserRepo) CreateUser(user *models.User) error {
	user.ID = uuid.New()
	return u.db.Create(user).Error
}

func (u *UserRepo) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := u.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
