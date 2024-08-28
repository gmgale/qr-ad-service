package db

import (
	"log"

	"github.com/gmgale/qr-ad-service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB initializes the database connection.
func ConnectDB(dsn string) {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Auto-migrate the schema.
	if err := DB.AutoMigrate(&models.User{}, &models.QRCode{}, &models.AdLog{}, &models.Revenue{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
