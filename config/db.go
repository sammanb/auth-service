package config

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

func loadEnv() error {
	viper.SetConfigName(".env") // file name without extension if you use SetConfigName
	viper.SetConfigType("env")  // tell Viper it's an env file
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

func PreInit() string {
	// dsn := "host=localhost user=postgres password=yourpass dbname=postgres sslmode=disable"
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable",
		viper.GetString("DB_HOST"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_PORT"),
	)
	adminDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("postgres db not found!")
	}

	dbName := viper.GetString("DB_NAME")
	var exists bool
	err = adminDB.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbName).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		_, err := adminDB.Exec("CREATE DATABASE " + dbName)
		if err != nil {
			log.Fatal("failed to create database. ", err)
		}
	}

	log.Println("database created successfully")
	return dbName
}

func InitDB() *gorm.DB {
	err := loadEnv()
	if err != nil {
		log.Fatal("failed to read env")
	}

	dbName := PreInit()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		viper.GetString("DB_HOST"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		dbName,
		viper.GetString("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	return db
}
