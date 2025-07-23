package utils

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	dsn := fmt.Sprintf("host=%s user=postgres password=postgres dbname=testdb port=5433 sslmode=disable", host)

	var db *gorm.DB
	var err error
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			return db
		}
		log.Println("Retrying test db connection...")
		time.Sleep(2 * time.Second)
	}

	log.Fatal("Could not connect to test db: ", err)
	// require.NoError(t, err)
	// require.NoError(t, db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}))
	// return db
	return nil
}
