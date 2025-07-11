package middleware

import (
	"log"
	"net/http"
	"os"
)

// SecureHTTPS redirects HTTP requests to HTTPS in production
func SecureHTTPS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only redirect in production environment
		if os.Getenv("ENVIRONMENT") == "production" &&
			r.Header.Get("X-Forwarded-Proto") == "http" {

			// Log the redirect for debugging
			log.Printf("Redirecting HTTP request to HTTPS: %s%s", r.Host, r.URL.Path)

			// Construct HTTPS URL
			target := "https://" + r.Host + r.URL.Path
			if r.URL.RawQuery != "" {
				target += "?" + r.URL.RawQuery
			}

			// Perform redirect
			http.Redirect(w, r, target, http.StatusTemporaryRedirect)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// SecurityHeaders adds security-related HTTP headers to responses
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only set security headers in production
		if os.Getenv("ENVIRONMENT") == "production" {
			// HSTS: Force browsers to use HTTPS for this domain
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

			// Prevent MIME type sniffing
			w.Header().Set("X-Content-Type-Options", "nosniff")

			// Prevent clickjacking
			w.Header().Set("X-Frame-Options", "DENY")

			// Control how much information is sent in the Referer header
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		}

		next.ServeHTTP(w, r)
	})
}
