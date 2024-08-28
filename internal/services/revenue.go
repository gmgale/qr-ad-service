package services

import (
	"fmt"
	"time"

	"github.com/gmgale/qr-ad-service/internal/db"
	"github.com/gmgale/qr-ad-service/internal/models"
)

// RevenueService handles the logic related to revenue calculation and distribution
type RevenueService struct {
}

// NewRevenueService initializes a new RevenueService
func NewRevenueService() *RevenueService {
	return &RevenueService{}
}

// CalculateAndDistributeRevenue calculates the revenue for an ad click and records it
func (s *RevenueService) CalculateAndDistributeRevenue(userID int64, qrCodeID int64, cpc float64) error {
	// Create a new revenue record
	revenue := models.Revenue{
		UserID:    userID,
		Amount:    cpc, // Use the actual CPC value
		CreatedAt: time.Now(),
	}

	// Save the revenue record to the database
	if err := db.DB.Create(&revenue).Error; err != nil {
		return fmt.Errorf("failed to record revenue: %v", err)
	}

	return nil
}
