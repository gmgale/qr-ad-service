package db

import (
	"fmt"
	"log"
	"os"

	"github.com/gmgale/qr-ad-service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB initializes the database connection using environment variables.
func ConnectDB() {
	var err error
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "youruser"),
		getEnv("DB_PASSWORD", "yourpassword"),
		getEnv("DB_NAME", "qrads"),
		getEnv("DB_PORT", "5432"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Auto-migrate the schema.
	if err := DB.AutoMigrate(&models.User{}, &models.QRCode{}, &models.AdLog{}, &models.Revenue{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

// getEnv reads an environment variable or returns a default value if not set.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
