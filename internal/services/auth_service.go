package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/samvibes/vexop/auth-service/internal/models"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	HashPassword(password string) (string, error)
	CompareHashAndPassword(password, hashed []byte) bool
	GenerateJWT(user *models.User) (string, error)
}

type AuthServiceImpl struct{}

func NewAuthService() AuthService {
	return &AuthServiceImpl{}
}

func (a *AuthServiceImpl) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (a *AuthServiceImpl) CompareHashAndPassword(password, hashed []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashed, password)
	return err == nil
}

func (a *AuthServiceImpl) GenerateJWT(user *models.User) (string, error) {
	userID := user.ID
	claims := jwt.MapClaims{
		"id":  userID.String(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	jwtSecret := []byte(viper.GetString("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
