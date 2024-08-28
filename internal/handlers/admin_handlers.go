package handlers

import (
	"encoding/json"
	"github.com/gmgale/qr-ad-service/internal/models"
	"net/http"

	"github.com/gmgale/qr-ad-service/internal/db"
)

type AdminStatsResponse struct {
	TotalUsers     int64   `json:"total_users"`
	TotalRevenue   float64 `json:"total_revenue"`
	TotalAdsServed int64   `json:"total_ads_served"`
	TotalClicks    int64   `json:"total_clicks"`
}

// Get system statistics for admin
func (s *Server) GetAdminStats(w http.ResponseWriter, r *http.Request) {
	var stats AdminStatsResponse

	if err := db.DB.Model(&models.User{}).Count(&stats.TotalUsers).Error; err != nil {
		internal.WriteAPIError(w, internal.NewAPIError(http.StatusInternalServerError, "Failed to retrieve user count"))
		return
	}

	if err := db.DB.Model(&models.AdLog{}).Where("is_clicked = ?", true).Count(&stats.TotalClicks).Error; err != nil {
		internal.WriteAPIError(w, internal.NewAPIError(http.StatusInternalServerError, "Failed to retrieve click count"))
		return
	}

	if err := db.DB.Model(&models.AdLog{}).Count(&stats.TotalAdsServed).Error; err != nil {
		internal.WriteAPIError(w, internal.NewAPIError(http.StatusInternalServerError, "Failed to retrieve ads served count"))
		return
	}

	if err := db.DB.Model(&models.Revenue{}).Select("SUM(amount)").Scan(&stats.TotalRevenue).Error; err != nil {
		internal.WriteAPIError(w, internal.NewAPIError(http.StatusInternalServerError, "Failed to retrieve revenue"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
