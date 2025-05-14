package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joshuagudgel/toasted-coffee/backend/internal/auth"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/middleware"
)

func TestJWTAuth(t *testing.T) {
	tests := []struct {
		name           string
		setupAuth      func(r *http.Request)
		expectedStatus int
	}{
		{
			name: "Valid token",
			setupAuth: func(r *http.Request) {
				// Generate a valid token
				token, _ := auth.GenerateToken(1, "admin")
				r.Header.Set("Authorization", "Bearer "+token)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Missing authorization header",
			setupAuth: func(r *http.Request) {
				// Don't set any auth header
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid authorization format",
			setupAuth: func(r *http.Request) {
				r.Header.Set("Authorization", "InvalidFormat")
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid token",
			setupAuth: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer invalidtoken123")
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Expired token",
			setupAuth: func(r *http.Request) {
				// You would need to generate an expired token here
				// For testing purposes, you might modify your auth package to accept a custom expiration
				// Or use a known expired token for testing
				r.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEsInJvbGUiOiJhZG1pbiIsImV4cCI6MTY0MTAxMjM0NSwiaWF0IjoxNjQxMDA4NzQ1LCJpc3MiOiJ0b2FzdGVkLWNvZmZlZS1jbyJ9.invalidSignature")
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test handler that will be wrapped by our middleware
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check if claims were added to context
				claims, ok := auth.ExtractClaimsFromContext(r.Context())
				if !ok && tc.expectedStatus == http.StatusOK {
					t.Error("Claims not found in context but request should be authorized")
				}
				if ok && tc.expectedStatus != http.StatusOK {
					t.Errorf("Claims found in context but request should not be authorized: %+v", claims)
				}

				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			})

			// Wrap the test handler with our JWT middleware
			handler := middleware.JWTAuth(testHandler)

			// Create test request
			req := httptest.NewRequest("GET", "/api/v1/protected", nil)
			tc.setupAuth(req)

			// Create response recorder
			w := httptest.NewRecorder()

			// Call the handler
			handler.ServeHTTP(w, req)

			// Check status
			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}
		})
	}
}
