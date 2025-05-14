package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/joshuagudgel/toasted-coffee/backend/internal/database"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/handlers"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/models"
)

// MockBookingRepository implements the repository interface for testing
type MockBookingRepository struct {
	CreateFunc    func(context.Context, *models.Booking) (int, error)
	CreateCalled  bool
	CreateBooking *models.Booking
	GetByIDFunc   func(context.Context, int) (*models.Booking, error)
	GetAllFunc    func(context.Context) ([]*models.Booking, error)
}

func (m *MockBookingRepository) Create(ctx context.Context, booking *models.Booking) (int, error) {
	m.CreateCalled = true
	m.CreateBooking = booking
	return m.CreateFunc(ctx, booking)
}

func (m *MockBookingRepository) GetByID(ctx context.Context, id int) (*models.Booking, error) {
	return m.GetByIDFunc(ctx, id)
}

func (m *MockBookingRepository) GetAll(ctx context.Context) ([]*models.Booking, error) {
	return m.GetAllFunc(ctx)
}

// Make sure MockBookingRepository implements the interface
var _ database.BookingRepositoryInterface = &MockBookingRepository{}

func TestCreateBookingHandler(t *testing.T) {
	tests := []struct {
		name           string
		booking        models.Booking
		mockCreateFunc func(context.Context, *models.Booking) (int, error)
		expectedStatus int
		expectedErr    string
	}{
		{
			name: "Valid booking with email",
			booking: models.Booking{
				Name:          "Test User",
				Email:         "test@example.com",
				Date:          "2025-06-01",
				Time:          "14:00",
				People:        5,
				Location:      "Test Location",
				CoffeeFlavors: []string{"french_toast"},
				MilkOptions:   []string{"whole"},
				Package:       "Group",
			},
			mockCreateFunc: func(ctx context.Context, b *models.Booking) (int, error) {
				return 1, nil
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Valid booking with phone",
			booking: models.Booking{
				Name:          "Test User",
				Phone:         "555-1234",
				Date:          "2025-06-01",
				Time:          "14:00",
				People:        5,
				Location:      "Test Location",
				CoffeeFlavors: []string{"french_toast"},
				MilkOptions:   []string{"whole"},
			},
			mockCreateFunc: func(ctx context.Context, b *models.Booking) (int, error) {
				return 2, nil
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Missing both email and phone",
			booking: models.Booking{
				Name:          "Test User",
				Date:          "2025-06-01",
				Time:          "14:00",
				People:        5,
				Location:      "Test Location",
				CoffeeFlavors: []string{"french_toast"},
				MilkOptions:   []string{"whole"},
			},
			mockCreateFunc: func(ctx context.Context, b *models.Booking) (int, error) {
				return 0, nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedErr:    "Email or phone number is required",
		},
		{
			name: "Malformed date",
			booking: models.Booking{
				Name:          "Test User",
				Email:         "test@example.com",
				Date:          "invalid-date",
				Time:          "14:00",
				People:        5,
				Location:      "Test Location",
				CoffeeFlavors: []string{"french_toast"},
				MilkOptions:   []string{"whole"},
			},
			mockCreateFunc: func(ctx context.Context, b *models.Booking) (int, error) {
				return 0, nil
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &MockBookingRepository{
				CreateFunc: func(ctx context.Context, booking *models.Booking) (int, error) {
					// Validate inputs if desired
					if booking.Name == "" {
						return 0, fmt.Errorf("name is required")
					}
					return 123, nil // Return ID 123 on success
				},
				GetByIDFunc: func(ctx context.Context, id int) (*models.Booking, error) {
					return &models.Booking{ID: id, Name: "Test"}, nil
				},
				GetAllFunc: func(ctx context.Context) ([]*models.Booking, error) {
					return []*models.Booking{{ID: 1, Name: "Test"}}, nil
				},
			}

			// Create handler with mock
			handler := handlers.NewBookingHandler(mockRepo)

			// Create request body
			body, _ := json.Marshal(tc.booking)
			req := httptest.NewRequest("POST", "/api/v1/bookings", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Call handler
			handler.Create(w, req)

			// Check status code
			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			// Check error message if expected
			if tc.expectedErr != "" {
				responseBody := w.Body.String()
				if !strings.Contains(responseBody, tc.expectedErr) {
					t.Errorf("Expected error '%s', got '%s'", tc.expectedErr, responseBody)
				}
			}
		})
	}
}

func TestGetAllBookingsHandler(t *testing.T) {
	tests := []struct {
		name           string
		mockGetAllFunc func(context.Context) ([]*models.Booking, error)
		expectedStatus int
		expectedCount  int
		expectedErr    string
	}{
		{
			name: "Successfully retrieve bookings",
			mockGetAllFunc: func(ctx context.Context) ([]*models.Booking, error) {
				return []*models.Booking{
					{ID: 1, Name: "User1"},
					{ID: 2, Name: "User2"},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name: "Empty bookings list",
			mockGetAllFunc: func(ctx context.Context) ([]*models.Booking, error) {
				return []*models.Booking{}, nil
			},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name: "Database error",
			mockGetAllFunc: func(ctx context.Context) ([]*models.Booking, error) {
				return nil, fmt.Errorf("database connection error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedErr:    "Failed to retrieve bookings",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &MockBookingRepository{
				GetAllFunc: tc.mockGetAllFunc,
			}

			// Create handler with mock
			handler := handlers.NewBookingHandler(mockRepo)

			// Create request
			req := httptest.NewRequest("GET", "/api/v1/bookings", nil)
			w := httptest.NewRecorder()

			// Call handler
			handler.GetAll(w, req)

			// Check status code
			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			// If successful, check the count of bookings
			if tc.expectedStatus == http.StatusOK {
				var bookings []*models.Booking
				if err := json.Unmarshal(w.Body.Bytes(), &bookings); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if len(bookings) != tc.expectedCount {
					t.Errorf("Expected %d bookings, got %d", tc.expectedCount, len(bookings))
				}
			}

			// Check error message if expected
			if tc.expectedErr != "" {
				responseBody := w.Body.String()
				if !strings.Contains(responseBody, tc.expectedErr) {
					t.Errorf("Expected error '%s', got '%s'", tc.expectedErr, responseBody)
				}
			}
		})
	}
}
