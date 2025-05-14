package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joshuagudgel/toasted-coffee/backend/internal/middleware"
)

func TestCORS(t *testing.T) {
	tests := []struct {
		name                 string
		allowedOrigins       string
		requestOrigin        string
		expectedAllowOrigin  string
		expectedAllowMethods bool
		expectedAllowHeaders bool
	}{
		{
			name:                 "Origin in allowed list",
			allowedOrigins:       "http://localhost:5173,https://toasted-coffee-frontend.onrender.com",
			requestOrigin:        "http://localhost:5173",
			expectedAllowOrigin:  "http://localhost:5173",
			expectedAllowMethods: true,
			expectedAllowHeaders: true,
		},
		{
			name:                 "Origin not in allowed list",
			allowedOrigins:       "http://localhost:5173,https://toasted-coffee-frontend.onrender.com",
			requestOrigin:        "http://evil-site.com",
			expectedAllowOrigin:  "",
			expectedAllowMethods: false,
			expectedAllowHeaders: false,
		},
		{
			name:                 "Space in allowed origins",
			allowedOrigins:       "http://localhost:5173, https://toasted-coffee-frontend.onrender.com",
			requestOrigin:        "https://toasted-coffee-frontend.onrender.com",
			expectedAllowOrigin:  "https://toasted-coffee-frontend.onrender.com",
			expectedAllowMethods: true,
			expectedAllowHeaders: true,
		},
		{
			name:                 "No origin in request",
			allowedOrigins:       "http://localhost:5173,https://toasted-coffee-frontend.onrender.com",
			requestOrigin:        "",
			expectedAllowOrigin:  "",
			expectedAllowMethods: false,
			expectedAllowHeaders: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test handler
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			})

			// Wrap with CORS middleware
			handler := middleware.CORS(tc.allowedOrigins)(testHandler)

			// Create test request
			req := httptest.NewRequest("GET", "/api/v1/bookings", nil)
			if tc.requestOrigin != "" {
				req.Header.Set("Origin", tc.requestOrigin)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Call the handler
			handler.ServeHTTP(w, req)

			// Check CORS headers
			allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
			if allowOrigin != tc.expectedAllowOrigin {
				t.Errorf("Expected Access-Control-Allow-Origin: %q, got %q", tc.expectedAllowOrigin, allowOrigin)
			}

			allowMethods := w.Header().Get("Access-Control-Allow-Methods")
			if (allowMethods != "") != tc.expectedAllowMethods {
				t.Errorf("Expected Access-Control-Allow-Methods to be %t, but was %q", tc.expectedAllowMethods, allowMethods)
			}

			allowHeaders := w.Header().Get("Access-Control-Allow-Headers")
			if (allowHeaders != "") != tc.expectedAllowHeaders {
				t.Errorf("Expected Access-Control-Allow-Headers to be %t, but was %q", tc.expectedAllowHeaders, allowHeaders)
			}
		})
	}
}

func TestCORS_OptionsRequest(t *testing.T) {
	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This should not be called for OPTIONS
		t.Error("Handler was called for OPTIONS request")
		w.WriteHeader(http.StatusTeapot) // Something unexpected
	})

	// Wrap with CORS middleware
	handler := middleware.CORS("http://localhost:5173")(testHandler)

	// Create OPTIONS request
	req := httptest.NewRequest("OPTIONS", "/api/v1/bookings", nil)
	req.Header.Set("Origin", "http://localhost:5173")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call the handler
	handler.ServeHTTP(w, req)

	// For OPTIONS request, should return 200 OK and not call the handler
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for OPTIONS, got %d", http.StatusOK, w.Code)
	}
}
