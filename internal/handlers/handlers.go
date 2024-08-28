package handlers

import (
	"net/http"

	"github.com/gmgale/qr-ad-service/internal/api"
	"github.com/gmgale/qr-ad-service/internal/auth"
	"github.com/go-chi/chi/v5"
)

// Implement the interface generated by oapi-codegen
type Server struct{}

// Get owner account summary
func (s *Server) GetOwnersSummary(w http.ResponseWriter, r *http.Request) {
	// Implementation here
}

// SetupRouter initializes the router and returns it
func SetupRouter() chi.Router {
	r := chi.NewRouter()
	s := &Server{}

	// Public routes
	r.Post("/auth/register", s.PostAuthRegister)
	r.Post("/auth/login", s.PostAuthLogin)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware) // Apply the JWT authentication middleware

		r.Get("/owners/summary", s.GetOwnersSummary)
		r.Post("/owners/qrcode", s.PostOwnersQrcode)
		r.Get("/admin/stats", s.GetAdminStats)
	})

	// Register the API routes with the server instance
	api.HandlerFromMux(s, r)

	return r
}
