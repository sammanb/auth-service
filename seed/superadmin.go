package seed

import (
	"log"

	"github.com/google/uuid"
	"github.com/samvibes/vexop/auth-service/internal/models"
	"github.com/samvibes/vexop/auth-service/internal/services"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func SeedSuperAdmin(db *gorm.DB) {
	email := viper.GetString("SUPERADMIN_EMAIL")
	password := viper.GetString("SUPERADMIN_PASSWORD")

	if email == "" || password == "" {
		log.Println("Superadmin credentials not provided. Skipping seed.")
		return
	}

	var saUser models.User
	err := db.Where("email = ?", email).First(&saUser).Error
	if err == nil {
		log.Println("Superadmin already exists. Skipping seed")
		return
	}

	hashed, err := services.HashPassword(password)
	if err != nil {
		log.Println("Superadmin password hash error. Skipping seed")
		return
	}

	superadmin := models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hashed,
		Role:         models.RoleSuperAdmin,
	}
	db.Create(&superadmin)
	log.Println("Superadmin created: ", email)
}
