package handlers

import (
	"encoding/json"
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

	db.DB.Model(&User{}).Count(&stats.TotalUsers)
	db.DB.Model(&AdLog{}).Where("is_clicked = ?", true).Count(&stats.TotalClicks)
	db.DB.Model(&AdLog{}).Count(&stats.TotalAdsServed)
	db.DB.Model(&Revenue{}).Select("SUM(amount)").Scan(&stats.TotalRevenue)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
