package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/database"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/handlers"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/models"
)

// MockBookingRepository implements the repository interface for testing
type MockBookingRepository struct {
	// Create
	CreateFunc    func(context.Context, *models.Booking) (int, error)
	CreateCalled  bool
	CreateBooking *models.Booking

	// GetByID
	GetByIDFunc   func(context.Context, int) (*models.Booking, error)
	GetByIDCalled bool
	GetByIDArg    int

	// GetAll
	GetAllFunc   func(context.Context) ([]*models.Booking, error)
	GetAllCalled bool
}

// Implement interface methods with tracking
func (m *MockBookingRepository) Create(ctx context.Context, booking *models.Booking) (int, error) {
	m.CreateCalled = true
	m.CreateBooking = booking
	return m.CreateFunc(ctx, booking)
}

func (m *MockBookingRepository) GetByID(ctx context.Context, id int) (*models.Booking, error) {
	m.GetByIDCalled = true
	m.GetByIDArg = id
	return m.GetByIDFunc(ctx, id)
}

func (m *MockBookingRepository) GetAll(ctx context.Context) ([]*models.Booking, error) {
	m.GetAllCalled = true
	return m.GetAllFunc(ctx)
}

// Verify interface implementation
var _ database.BookingRepositoryInterface = &MockBookingRepository{}

func TestCreateBookingHandler(t *testing.T) {
	tests := []struct {
		name           string
		booking        models.Booking
		mockCreateFunc func(context.Context, *models.Booking) (int, error)
		expectedStatus int
		expectedErr    string
		expectedID     int
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
				return 123, nil
			},
			expectedStatus: http.StatusCreated,
			expectedID:     123,
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
				return 124, nil
			},
			expectedStatus: http.StatusCreated,
			expectedID:     124,
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
		{
			name: "Database error",
			booking: models.Booking{
				Name:          "Test User",
				Email:         "test@example.com",
				Date:          "2025-06-01",
				Time:          "14:00",
				People:        5,
				Location:      "Test Location",
				CoffeeFlavors: []string{"french_toast"},
				MilkOptions:   []string{"whole"},
			},
			mockCreateFunc: func(ctx context.Context, b *models.Booking) (int, error) {
				return 0, fmt.Errorf("database connection error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedErr:    "Failed to create booking",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock repository with the test case's function
			mockRepo := &MockBookingRepository{
				CreateFunc: tc.mockCreateFunc,
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

			// Check if the mock was called when expected
			if tc.expectedStatus == http.StatusCreated && !mockRepo.CreateCalled {
				t.Error("Expected Create method to be called, but it wasn't")
			}

			// Check for success response
			if tc.expectedStatus == http.StatusCreated {
				var resp map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}

				// Check ID in response
				id, ok := resp["id"].(float64)
				if !ok {
					t.Error("Expected 'id' field in response")
				} else if int(id) != tc.expectedID {
					t.Errorf("Expected ID %d, got %d", tc.expectedID, int(id))
				}

				// Verify message is present
				if _, ok := resp["message"]; !ok {
					t.Error("Expected 'message' field in response")
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

				// For GetAll empty response test
				if tc.expectedCount == 0 {
					// Should still be a valid JSON array
					if w.Body.String() != "[]" && w.Body.String() != "[]\n" {
						t.Errorf("Expected empty JSON array, got: %s", w.Body.String())
					}
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

func TestGetBookingByIDHandler(t *testing.T) {
	tests := []struct {
		name            string
		bookingID       string
		mockGetByIDFunc func(context.Context, int) (*models.Booking, error)
		expectedStatus  int
		expectedErr     string
	}{
		{
			name:      "Valid booking ID",
			bookingID: "123",
			mockGetByIDFunc: func(ctx context.Context, id int) (*models.Booking, error) {
				return &models.Booking{
					ID:            id,
					Name:          "Test User",
					Email:         "test@example.com",
					Date:          "2025-06-01",
					Time:          "14:00",
					People:        5,
					Location:      "Test Location",
					CoffeeFlavors: []string{"french_toast"},
					MilkOptions:   []string{"whole"},
					Package:       "Group",
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "Non-existent booking ID",
			bookingID: "999",
			mockGetByIDFunc: func(ctx context.Context, id int) (*models.Booking, error) {
				return nil, fmt.Errorf("booking not found")
			},
			expectedStatus: http.StatusNotFound,
			expectedErr:    "Booking not found",
		},
		{
			name:      "Invalid booking ID format",
			bookingID: "abc",
			mockGetByIDFunc: func(ctx context.Context, id int) (*models.Booking, error) {
				return nil, nil // Should not be called
			},
			expectedStatus: http.StatusBadRequest,
			expectedErr:    "Invalid booking ID",
		},
		{
			name:      "Database error",
			bookingID: "123",
			mockGetByIDFunc: func(ctx context.Context, id int) (*models.Booking, error) {
				return nil, fmt.Errorf("database error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedErr:    "Failed to retrieve booking",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &MockBookingRepository{
				GetByIDFunc: tc.mockGetByIDFunc,
			}

			// Create handler with mock
			handler := handlers.NewBookingHandler(mockRepo)

			// Create request with URL parameter
			req := httptest.NewRequest("GET", "/api/v1/bookings/"+tc.bookingID, nil)

			// Setup chi context with URL parameters
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.bookingID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Create response recorder
			w := httptest.NewRecorder()

			// Call handler
			handler.GetByID(w, req)

			// Check status code
			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			// For valid ID, check that the response contains booking data
			if tc.expectedStatus == http.StatusOK {
				var booking models.Booking
				if err := json.Unmarshal(w.Body.Bytes(), &booking); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}

				// Verify ID was passed to repository
				id, _ := strconv.Atoi(tc.bookingID)
				if mockRepo.GetByIDArg != id {
					t.Errorf("Expected GetByID called with %d, got %d", id, mockRepo.GetByIDArg)
				}

				// Verify booking properties
				if booking.ID != id {
					t.Errorf("Expected booking ID %d, got %d", id, booking.ID)
				}

				if booking.Name == "" {
					t.Error("Expected non-empty booking name")
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

func TestResponseHeaders(t *testing.T) {
	// Create mock repository
	mockRepo := &MockBookingRepository{
		GetAllFunc: func(ctx context.Context) ([]*models.Booking, error) {
			return []*models.Booking{}, nil
		},
	}

	// Create handler with mock
	handler := handlers.NewBookingHandler(mockRepo)

	// Create request
	req := httptest.NewRequest("GET", "/api/v1/bookings", nil)

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	handler.GetAll(w, req)

	// Check Content-Type header
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	// Check response is valid JSON
	var response interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Response is not valid JSON: %v", err)
	}
}
