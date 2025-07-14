package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/samvibes/vexop/auth-service/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CompareHashAndPassword(password, hashed []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashed, password)
	return err == nil
}

var jwtSecret = []byte("secret-key")

func GenerateJWT(user *models.User) (string, error) {
	userID := user.ID
	tenantID := user.TenantID
	role := user.Role
	claims := jwt.MapClaims{
		"id":   userID.String(),
		"role": role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	}

	if tenantID != nil {
		claims["tenant_id"] = tenantID.String()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
