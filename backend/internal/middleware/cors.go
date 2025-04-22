package middleware

import (
	"net/http"
	"strings"
)

// CORS adds CORS headers to responses
func CORS(allowOrigins string) func(next http.Handler) http.Handler {
	origins := strings.Split(allowOrigins, ",")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if the request origin is in our allowed list
			allowOrigin := ""
			for _, allowed := range origins {
				if origin == strings.TrimSpace(allowed) {
					allowOrigin = origin
					break
				}
			}

			// set header only if origin is allowed
			if allowOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
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
