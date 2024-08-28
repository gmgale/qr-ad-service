package services

import (
	"fmt"
	"log"

	"github.com/gmgale/qr-ad-service/internal/db"
	"github.com/gmgale/qr-ad-service/internal/models"
	"github.com/kritzware/google-ads-go/ads"
	"github.com/kritzware/google-ads-go/services"
)

// GoogleAdsService wraps the Google Ads API client
type GoogleAdsService struct {
	client *ads.GoogleAdsClient
}

// NewGoogleAdsService initializes the Google Ads client using credentials from a storage file
func NewGoogleAdsService(credentialsPath string) (*GoogleAdsService, error) {
	client, err := ads.NewClientFromStorage(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create Google Ads client: %v", err)
	}

	return &GoogleAdsService{client: client}, nil
}

// ServeAd serves an ad and tracks its view
func (s *GoogleAdsService) ServeAd(qrCodeID int64) error {
	// Load the GoogleAds service
	googleAdsService := services.NewGoogleAdsServiceClient(s.client.Conn())

	// Create a search request to retrieve campaign data
	request := services.SearchGoogleAdsRequest{
		CustomerId: "YOUR_CUSTOMER_ID", // Replace with your customer ID
		Query:      "SELECT campaign.id, ad_group_ad.ad.id, metrics.average_cpc FROM ad_group_ad WHERE campaign.status = 'ENABLED' AND ad_group_ad.status = 'ENABLED'",
	}

	// Get the results
	response, err := googleAdsService.Search(s.client.Context(), &request)
	if err != nil {
		return fmt.Errorf("failed to serve ad: %v", err)
	}

	// Process the response to find the CPC
	var averageCPC float64
	for _, row := range response.Results {
		ad := row.GetAdGroupAd().GetAd()
		averageCPC = row.GetMetrics().GetAverageCpc().Value / 1_000_000 // Converting from micros to base currency
		log.Printf("Served Ad ID: %d with CPC: %.2f", ad.Id.Value, averageCPC)
	}

	log.Println("Ad served for QR Code ID:", qrCodeID)
	return nil
}

// TrackAdClick tracks a click on the served ad
func (s *GoogleAdsService) TrackAdClick(qrCodeID int64) error {
	// Retrieve the QR code to find the associated owner
	var qrCode models.QRCode
	if err := db.DB.First(&qrCode, qrCodeID).Error; err != nil {
		return fmt.Errorf("failed to retrieve QR code: %v", err)
	}

	// Calculate and distribute revenue using the actual CPC
	averageCPC, err := s.GetAverageCPC(qrCodeID)
	if err != nil {
		return fmt.Errorf("failed to get average CPC: %v", err)
	}

	revenueService := NewRevenueService()
	if err := revenueService.CalculateAndDistributeRevenue(qrCode.UserID, qrCodeID, averageCPC); err != nil {
		return fmt.Errorf("failed to calculate and distribute revenue: %v", err)
	}

	log.Printf("Ad click tracked and revenue distributed for QR Code ID: %d", qrCodeID)
	return nil
}

// GetAverageCPC retrieves the average CPC for a given ad or campaign
func (s *GoogleAdsService) GetAverageCPC(qrCodeID int64) (float64, error) {
	// Load the GoogleAds service
	googleAdsService := services.NewGoogleAdsServiceClient(s.client.Conn())

	// Create a search request to retrieve the average CPC
	request := services.SearchGoogleAdsRequest{
		CustomerId: "YOUR_CUSTOMER_ID",                                                                 // Replace with your customer ID
		Query:      "SELECT metrics.average_cpc FROM ad_group_ad WHERE ad_group_ad.ad.id = YOUR_AD_ID", // Adjust the query to target specific ad or campaign
	}

	// Get the results
	response, err := googleAdsService.Search(s.client.Context(), &request)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve average CPC: %v", err)
	}

	// Extract and return the average CPC
	for _, row := range response.Results {
		averageCPC := row.GetMetrics().GetAverageCpc().Value / 1_000_000 // Converting from micros to base currency
		return averageCPC, nil
	}

	return 0, fmt.Errorf("no CPC data found")
}
