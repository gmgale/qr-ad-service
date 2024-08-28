package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// UserContextKey is the key used to store the user ID in the context
type UserContextKey string

const UserIDContextKey = UserContextKey("user_id")

// AuthMiddleware is a middleware that protects routes by validating the JWT token
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := VerifyJWT(tokenString)
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
			} else {
				http.Error(w, "Failed to parse token", http.StatusBadRequest)
			}
			return
		}

		// Store the user ID in the context
		ctx := context.WithValue(r.Context(), UserIDContextKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
