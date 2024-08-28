package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gmgale/qr-ad-service/internal/db"
	"github.com/gmgale/qr-ad-service/internal/models"
	"github.com/skip2/go-qrcode"
)

// Generate a new QR code
func (s *Server) PostOwnersQrcode(w http.ResponseWriter, r *http.Request) {
	var req models.QRCode
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Generate the QR code
	png, err := qrcode.Encode(req.OriginalURL, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		return
	}

	req.GeneratedCode = string(png)
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	if err := db.DB.Create(&req).Error; err != nil {
		http.Error(w, "Failed to save QR code", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(png) // Optionally return the PNG data directly
}
