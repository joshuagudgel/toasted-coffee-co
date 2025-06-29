package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"
	"log"

	"github.com/joshuagudgel/toasted-coffee/backend/internal/auth"
)

// JWTAuth middleware intercepts requests to validate JWT tokens
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		log.Printf("JWT VALIDATION START: Request to %s", r.URL.Path)
		
		// First try to get token from cookie
		tokenCookie, err := r.Cookie("access_token")
		var tokenString string
		
		if err == nil {
			// Found token in cookie
			tokenString = tokenCookie.Value
		} else {
			// Fall back to Authorization header for API clients
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Printf("JWT VALIDATION: No token found for %s", r.URL.Path)
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			// Extract token from Bearer scheme
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				log.Printf("JWT VALIDATION: Invalid authorization format for %s", r.URL.Path)
				http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
				return
			}
			
			tokenString = tokenParts[1]
		}

		// Validate token using the auth package
		validateStart := time.Now()
		claims, err := auth.ValidateToken(tokenString)
		validationTime := time.Since(validateStart)
		log.Printf("JWT VALIDATION TIMING: Token validation took %v for %s", validationTime, r.URL.Path)
		
		if err != nil {
			log.Printf("JWT VALIDATION: Invalid token for %s: %v", r.URL.Path, err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add claims to context using the exported key from auth
		ctx := context.WithValue(r.Context(), auth.ClaimsContextKey, claims)
		
		totalTime := time.Since(startTime)
		log.Printf("JWT VALIDATION COMPLETE: Total processing time %v for %s", totalTime, r.URL.Path)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
