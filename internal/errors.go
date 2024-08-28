package internal

import (
	"fmt"
	"net/http"
)

// APIError defines a structured error for the API
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface for APIError
func (e *APIError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

// NewAPIError creates a new APIError with a specific HTTP status code and message
func NewAPIError(code int, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

// WriteAPIError writes the APIError as a JSON response
func WriteAPIError(w http.ResponseWriter, err *APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)
	fmt.Fprintf(w, `{"error": "%s"}`, err.Message)
}
