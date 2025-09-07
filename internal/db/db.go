package db

import (
	"log"
	"os"

	"github.com/adriel-meb/appointly-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// DbConnect connects to the database using environment variables
func DbConnect() {
	var err error
	dsn := "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_NAME") + " port=" + os.Getenv("DB_PORT") + " sslmode=disable TimeZone=Africa/Libreville"

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect database. DSN = ", dsn, " error: ", err)
	}

	log.Println("Database connected")
}

func DbMigration() {
	DbConnect()
	err := DB.AutoMigrate(
		&models.User{}, &models.Provider{}, &models.Specialization{},
	)

	if err != nil {
		log.Fatal("failed to migrate database: ", err)
	}

	log.Println("Database migrated")
}
