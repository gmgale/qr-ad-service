package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gmgale/qr-ad-service/internal/db"
	"github.com/gmgale/qr-ad-service/internal/handlers"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Database connection string (you can load this from environment variables)
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=youruser password=yourpassword dbname=qrads port=5432 sslmode=disable"
	}

	// Initialize the database connection
	db.ConnectDB(dsn)

	// Set up the router
	r := handlers.SetupRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
