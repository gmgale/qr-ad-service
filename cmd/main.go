package main

import (
	"log"
	"net/http"

	"github.com/gmgale/qr-ad-service/internal/handlers"
)

func main() {
	r := handlers.SetupRouter()

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
