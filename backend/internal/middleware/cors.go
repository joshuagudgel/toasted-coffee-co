package middleware

import (
	"net/http"
	"os"
	"strings"
	"log"
)

// CORS adds CORS headers to responses
func CORS(allowOrigins string) func(next http.Handler) http.Handler {
	origins := strings.Split(allowOrigins, ",")

	// Force HTTPS for non-localhost origins in production
	if os.Getenv("ENVIRONMENT") == "production" {
		for i, origin := range origins {
			origin = strings.TrimSpace(origin)
			if !strings.Contains(origin, "localhost") && strings.HasPrefix(origin, "http:") {
				origins[i] = "https:" + strings.TrimPrefix(origin, "http:")
				log.Printf("Converted origin from HTTP to HTTPS: %s -> %s", origin, origins[i])
			} else {
				origins[i] = origin
			}
		}
	} else {
		// Just trim whitespace in development
		for i, origin := range origins {
			origins[i] = strings.TrimSpace(origin)
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if the request origin is in our allowed list
			allowOrigin := ""
			for _, allowed := range origins {
				if origin == allowed {
					allowOrigin = origin // Use the exact origin
					break
				}
			}

			// Set headers only if origin is allowed
			if allowOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Max-Age", "3600")
			}

			// handle preflight OPTIONS request
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
