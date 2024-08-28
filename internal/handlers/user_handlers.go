package handlers

import (
	"encoding/json"
	"github.com/gmgale/qr-ad-service/internal"
	"github.com/gmgale/qr-ad-service/internal/auth"
	"net/http"
	"time"

	"github.com/gmgale/qr-ad-service/internal/db"
	"github.com/gmgale/qr-ad-service/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// Register a new owner account
func (s *Server) PostAuthRegister(w http.ResponseWriter, r *http.Request) {
	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		internal.WriteAPIError(w, internal.NewAPIError(http.StatusBadRequest, "Invalid request payload"))
		return
	}

	// Check if the email already exists
	var existingUser models.User
	if err := db.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		internal.WriteAPIError(w, internal.NewAPIError(http.StatusConflict, "Email already in use"))
		return
	}

	// Hash the password before saving it to the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		internal.WriteAPIError(w, internal.NewAPIError(http.StatusInternalServerError, "Failed to hash password"))
		return
	}

	// Create the new user record
	newUser := models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		BusinessName: req.BusinessName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := db.DB.Create(&newUser).Error; err != nil {
		internal.WriteAPIError(w, internal.NewAPIError(http.StatusInternalServerError, "Failed to create user"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// Login an owner
func (s *Server) PostAuthLogin(w http.ResponseWriter, r *http.Request) {
	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		internal.WriteAPIError(w, internal.NewAPIError(http.StatusBadRequest, "Invalid request payload"))
		return
	}

	var user models.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		internal.WriteAPIError(w, internal.NewAPIError(http.StatusUnauthorized, "Invalid credentials"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.PasswordHash)); err != nil {
		internal.WriteAPIError(w, internal.NewAPIError(http.StatusUnauthorized, "Invalid credentials"))
		return
	}

	// Generate JWT token
	tokenString, err := auth.GenerateJWT(user.ID)
	if err != nil {
		internal.WriteAPIError(w, internal.NewAPIError(http.StatusInternalServerError, "Failed to generate token"))
		return
	}

	// Return the token as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
