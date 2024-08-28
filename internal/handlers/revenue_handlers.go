package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gmgale/qr-ad-service/internal/auth"
	"github.com/gmgale/qr-ad-service/internal/db"
	"github.com/gmgale/qr-ad-service/internal/models"
)

// GetOwnerRevenue returns the total revenue for the logged-in owner
func (s *Server) GetOwnerRevenue(w http.ResponseWriter, r *http.Request) {
	// Extract the user ID from the JWT context
	userID := r.Context().Value(auth.UserIDContextKey).(int64)

	var revenueTotal float64
	err := db.DB.Model(&models.Revenue{}).Where("user_id = ?", userID).Select("SUM(amount)").Scan(&revenueTotal).Error
	if err != nil {
		http.Error(w, "Failed to retrieve revenue", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"total_revenue": revenueTotal})
}
