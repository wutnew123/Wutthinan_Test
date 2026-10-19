package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *sql.DB
var GormDB *gorm.DB

func ConnectDB() {
	// Load values from .env when running locally. In production, environment
	// variables can be provided by the host/container instead.
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found; using system environment variables")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)
	var err error

	GormDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Cannot initialize PostgreSQL connection:", err)
	}

	DB, err = GormDB.DB()
	if err != nil {
		log.Fatal("Cannot get database connection:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("DB ping failed:", err)
	}

	log.Println("Connected to PostgreSQL")
}
