package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/database"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/handlers"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/middleware"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/models"
)

func setupTestServer(t *testing.T) (*chi.Mux, *database.DB, func()) {
	// Get test database URL
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/toasted_coffee_test?sslmode=disable"
	}

	// Connect to database
	db, err := database.New(dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Create tables
	_, err = db.Pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS bookings (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            email VARCHAR(255),
            phone VARCHAR(20),
            date DATE NOT NULL,
            time VARCHAR(10) NOT NULL,
            people INTEGER NOT NULL,
            location VARCHAR(255) NOT NULL,
            notes TEXT,
            coffee_flavors VARCHAR[] NOT NULL,
            milk_options VARCHAR[] NOT NULL,
            package VARCHAR(100),
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		t.Fatalf("Failed to create test tables: %v", err)
	}

	// Create repositories
	bookingRepo := database.NewBookingRepository(db)

	// Create handlers
	bookingHandler := handlers.NewBookingHandler(bookingRepo)

	// Create router
	r := chi.NewRouter()

	// Add middlewares
	r.Use(middleware.CORS("*"))

	// Add routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/bookings", bookingHandler.Create)
		r.Get("/bookings", bookingHandler.GetAll)
		r.Get("/bookings/{id}", bookingHandler.GetByID)
	})

	// Return cleanup function
	cleanup := func() {
		_, err := db.Pool.Exec(context.Background(), "DELETE FROM bookings")
		if err != nil {
			t.Logf("Warning: Failed to clean up test data: %v", err)
		}
		db.Close()
	}

	return r, db, cleanup
}

func TestFullBookingFlow(t *testing.T) {
	// Skip integration tests if flag is set
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration tests")
	}

	// Setup test server
	r, _, cleanup := setupTestServer(t)
	defer cleanup()

	// Test data
	booking := models.Booking{
		Name:          "Integration Test",
		Email:         "integration@test.com",
		Date:          "2025-07-01",
		Time:          "15:00",
		People:        10,
		Location:      "Integration Test Location",
		CoffeeFlavors: []string{"french_toast", "mexican_mocha"},
		MilkOptions:   []string{"whole", "oat"},
		Package:       "Crowd",
	}

	// Step 1: Create booking
	body, _ := json.Marshal(booking)
	req := httptest.NewRequest("POST", "/api/v1/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status %d, got %d: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	var createResp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &createResp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Extract booking ID
	bookingID, ok := createResp["id"].(float64)
	if !ok {
		t.Fatalf("Expected numeric booking ID, got %T: %v", createResp["id"], createResp["id"])
	}

	// Step 2: Get booking by ID
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/bookings/%d", int(bookingID)), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var retrievedBooking models.Booking
	if err := json.Unmarshal(w.Body.Bytes(), &retrievedBooking); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Verify booking data
	if retrievedBooking.Name != booking.Name {
		t.Errorf("Expected name %s, got %s", booking.Name, retrievedBooking.Name)
	}

	if retrievedBooking.Email != booking.Email {
		t.Errorf("Expected email %s, got %s", booking.Email, retrievedBooking.Email)
	}

	if retrievedBooking.Package != booking.Package {
		t.Errorf("Expected package %s, got %s", booking.Package, retrievedBooking.Package)
	}

	// Step 3: Get all bookings
	req = httptest.NewRequest("GET", "/api/v1/bookings", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var bookings []models.Booking
	if err := json.Unmarshal(w.Body.Bytes(), &bookings); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(bookings) < 1 {
		t.Errorf("Expected at least 1 booking, got %d", len(bookings))
	}

	var found bool
	for _, b := range bookings {
		if b.ID == int(bookingID) {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Booking ID %d not found in bookings list", int(bookingID))
	}
}
