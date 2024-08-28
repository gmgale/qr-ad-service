package handlers

import (
	"encoding/json"
	"github.com/gmgale/qr-ad-service/internal"
	"github.com/gmgale/qr-ad-service/internal/auth"
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
		internal.WriteAPIError(w, internal.NewAPIError(http.StatusBadRequest, "Invalid request payload"))
		return
	}

	// Generate the QR code
	qrCodePNG, err := qrcode.Encode(req.OriginalURL, qrcode.Medium, 256)
	if err != nil {
		internal.WriteAPIError(w, internal.NewAPIError(http.StatusInternalServerError, "Failed to generate QR code"))
		return
	}

	// Extract the user ID from the JWT context
	userID := r.Context().Value(auth.UserIDContextKey).(int64)

	// Save the QR code data to the database
	newQRCode := models.QRCode{
		OriginalURL:   req.OriginalURL,
		GeneratedCode: string(qrCodePNG),
		UserID:        userID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := db.DB.Create(&newQRCode).Error; err != nil {
		internal.WriteAPIError(w, internal.NewAPIError(http.StatusInternalServerError, "Failed to save QR code"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(qrCodePNG) // Optionally return the PNG data directly
}
