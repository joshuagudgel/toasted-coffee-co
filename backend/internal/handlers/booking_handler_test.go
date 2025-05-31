package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	GetAllFunc            func(context.Context, bool) ([]*models.Booking, error)
	GetAllCalled          bool
	GetAllIncludeArchived bool

	// Delete
	DeleteFunc   func(context.Context, int) error
	DeleteCalled bool
	DeleteArg    int

	// Update
	UpdateFunc    func(context.Context, int, *models.Booking) error
	UpdateCalled  bool
	UpdateID      int
	UpdateBooking *models.Booking

	// Archive
	ArchiveFunc   func(context.Context, int) error
	ArchiveCalled bool
	ArchiveArg    int

	// Unarchive
	UnarchiveFunc   func(context.Context, int) error
	UnarchiveCalled bool
	UnarchiveArg    int
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

func (m *MockBookingRepository) GetAll(ctx context.Context, includeArchived bool) ([]*models.Booking, error) {
	m.GetAllCalled = true
	m.GetAllIncludeArchived = includeArchived
	return m.GetAllFunc(ctx, includeArchived)
}

func (m *MockBookingRepository) Delete(ctx context.Context, id int) error {
	m.DeleteCalled = true
	m.DeleteArg = id
	return m.DeleteFunc(ctx, id)
}
func (m *MockBookingRepository) Update(ctx context.Context, id int, booking *models.Booking) error {
	m.UpdateCalled = true
	m.UpdateID = id
	m.UpdateBooking = booking
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, booking)
	}
	return nil
}

func (m *MockBookingRepository) Archive(ctx context.Context, id int) error {
	m.ArchiveCalled = true
	m.ArchiveArg = id
	if m.ArchiveFunc != nil {
		return m.ArchiveFunc(ctx, id)
	}
	return nil
}

func (m *MockBookingRepository) Unarchive(ctx context.Context, id int) error {
	m.UnarchiveCalled = true
	m.UnarchiveArg = id
	if m.UnarchiveFunc != nil {
		return m.UnarchiveFunc(ctx, id)
	}
	return nil
}

// Verify interface implementation
var _ database.BookingRepositoryInterface = &MockBookingRepository{}

func TestCreateBookingHandler(t *testing.T) {
	log.Println("Starting TestCreateBookingHandler")
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

func TestUpdateBookingHandler(t *testing.T) {
	log.Println("Starting TestUpdateBookingHandler")
	tests := []struct {
		name            string
		bookingID       string
		updatedBooking  models.Booking
		mockGetByIDFunc func(context.Context, int) (*models.Booking, error)
		mockUpdateFunc  func(context.Context, int, *models.Booking) error
		expectedStatus  int
		expectedErr     string
	}{
		{
			name:      "Successfully update booking",
			bookingID: "123",
			updatedBooking: models.Booking{
				Name:          "Updated User",
				Email:         "updated@example.com",
				Date:          "2025-07-01",
				Time:          "15:00",
				People:        7,
				Location:      "Updated Location",
				CoffeeFlavors: []string{"vanilla_bean"},
				MilkOptions:   []string{"oat"},
				Package:       "Premium",
			},
			mockGetByIDFunc: func(ctx context.Context, id int) (*models.Booking, error) {
				return &models.Booking{
					ID:            id,
					Name:          "Original User",
					Email:         "original@example.com",
					Date:          "2025-06-01",
					Time:          "14:00",
					People:        5,
					Location:      "Original Location",
					CoffeeFlavors: []string{"french_toast"},
					MilkOptions:   []string{"whole"},
					Package:       "Standard",
				}, nil
			},
			mockUpdateFunc: func(ctx context.Context, id int, booking *models.Booking) error {
				return nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "Invalid booking ID format",
			bookingID: "abc",
			updatedBooking: models.Booking{
				Name:  "Updated User",
				Email: "updated@example.com",
			},
			mockGetByIDFunc: func(ctx context.Context, id int) (*models.Booking, error) {
				return nil, nil // Should not be called
			},
			mockUpdateFunc: func(ctx context.Context, id int, booking *models.Booking) error {
				return nil // Should not be called
			},
			expectedStatus: http.StatusBadRequest,
			expectedErr:    "Invalid booking ID",
		},
		{
			name:      "Booking not found",
			bookingID: "999",
			updatedBooking: models.Booking{
				Name:  "Updated User",
				Email: "updated@example.com",
			},
			mockGetByIDFunc: func(ctx context.Context, id int) (*models.Booking, error) {
				return nil, nil
			},
			mockUpdateFunc: func(ctx context.Context, id int, booking *models.Booking) error {
				return nil // Should not be called
			},
			expectedStatus: http.StatusNotFound,
			expectedErr:    "Booking not found",
		},
		{
			name:      "Invalid updated booking data",
			bookingID: "123",
			updatedBooking: models.Booking{
				// Missing required fields
				Name:          "",
				Email:         "",
				Phone:         "",
				Date:          "2025-07-01",
				Time:          "15:00",
				People:        0, // Invalid: must be > 0
				Location:      "Updated Location",
				CoffeeFlavors: []string{}, // Invalid: empty array
				MilkOptions:   []string{}, // Invalid: empty array
			},
			mockGetByIDFunc: func(ctx context.Context, id int) (*models.Booking, error) {
				return &models.Booking{ID: id, Name: "Original User"}, nil
			},
			mockUpdateFunc: func(ctx context.Context, id int, booking *models.Booking) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedErr:    "Email or phone number is required",
		},
		{
			name:      "Database error",
			bookingID: "123",
			updatedBooking: models.Booking{
				Name:          "Updated User",
				Email:         "updated@example.com",
				Date:          "2025-07-01",
				Time:          "15:00",
				People:        7,
				Location:      "Updated Location",
				CoffeeFlavors: []string{"vanilla_bean"},
				MilkOptions:   []string{"oat"},
			},
			mockGetByIDFunc: func(ctx context.Context, id int) (*models.Booking, error) {
				return &models.Booking{ID: id, Name: "Original User"}, nil
			},
			mockUpdateFunc: func(ctx context.Context, id int, booking *models.Booking) error {
				return fmt.Errorf("database error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedErr:    "Failed to update booking",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &MockBookingRepository{
				GetByIDFunc: tc.mockGetByIDFunc,
				UpdateFunc:  tc.mockUpdateFunc,
			}

			// Create handler with mock
			handler := handlers.NewBookingHandler(mockRepo)

			// Create request with URL parameter and body
			body, _ := json.Marshal(tc.updatedBooking)
			req := httptest.NewRequest("PUT", "/api/v1/bookings/"+tc.bookingID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Setup chi context with URL parameters
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.bookingID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Create response recorder
			w := httptest.NewRecorder()

			// Call handler
			handler.Update(w, req)

			// Check status code
			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			// For successful updates, verify repository methods were called correctly
			if tc.expectedStatus == http.StatusOK {
				id, _ := strconv.Atoi(tc.bookingID)

				// Verify GetByID was called
				if !mockRepo.GetByIDCalled {
					t.Error("Expected GetByID to be called, but it wasn't")
				}
				if mockRepo.GetByIDArg != id {
					t.Errorf("GetByID called with wrong ID, expected %d, got %d", id, mockRepo.GetByIDArg)
				}

				// Verify Update was called
				if !mockRepo.UpdateCalled {
					t.Error("Expected Update to be called, but it wasn't")
				}
				if mockRepo.UpdateID != id {
					t.Errorf("Update called with wrong ID, expected %d, got %d", id, mockRepo.UpdateID)
				}

				// Verify the booking passed to Update contains the updates
				if mockRepo.UpdateBooking != nil {
					updatedBooking := mockRepo.UpdateBooking
					if updatedBooking.Name != tc.updatedBooking.Name {
						t.Errorf("Expected updated name %s, got %s", tc.updatedBooking.Name, updatedBooking.Name)
					}
					if updatedBooking.Email != tc.updatedBooking.Email {
						t.Errorf("Expected updated email %s, got %s", tc.updatedBooking.Email, updatedBooking.Email)
					}
				}

				// Verify response contains success message
				var resp map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}

				if message, ok := resp["message"]; !ok || !strings.Contains(message.(string), "updated") {
					t.Errorf("Expected success message containing 'updated', got %v", message)
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
	log.Println("Starting TestGetAllBookingsHandler")
	tests := []struct {
		name           string
		mockGetAllFunc func(context.Context, bool) ([]*models.Booking, error)
		expectedStatus int
		expectedCount  int
		expectedErr    string
	}{
		{
			name: "Successfully retrieve bookings",
			mockGetAllFunc: func(ctx context.Context, includeArchived bool) ([]*models.Booking, error) {
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
			mockGetAllFunc: func(ctx context.Context, includeArchived bool) ([]*models.Booking, error) {
				return []*models.Booking{}, nil
			},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name: "Database error",
			mockGetAllFunc: func(ctx context.Context, includeArchived bool) ([]*models.Booking, error) {
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
	log.Println("Starting TestGetBookingByIDHandler")
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

func TestDeleteBookingHandler(t *testing.T) {
	log.Println("Starting TestDeleteBookingHandler")
	tests := []struct {
		name            string
		bookingID       string
		mockGetByIDFunc func(context.Context, int) (*models.Booking, error)
		mockDeleteFunc  func(context.Context, int) error
		expectedStatus  int
		expectedErr     string
	}{
		{
			name:      "Valid booking ID",
			bookingID: "123",
			mockGetByIDFunc: func(ctx context.Context, id int) (*models.Booking, error) {
				return &models.Booking{ID: id, Name: "Test User"}, nil
			},
			mockDeleteFunc: func(ctx context.Context, id int) error {
				return nil
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:      "Non-existent booking ID",
			bookingID: "999",
			mockGetByIDFunc: func(ctx context.Context, id int) (*models.Booking, error) {
				return nil, nil
			},
			mockDeleteFunc: func(ctx context.Context, id int) error {
				return nil // Should not be called
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
			mockDeleteFunc: func(ctx context.Context, id int) error {
				return nil // Should not be called
			},
			expectedStatus: http.StatusBadRequest,
			expectedErr:    "Invalid booking ID",
		},
		{
			name:      "Database error on delete",
			bookingID: "123",
			mockGetByIDFunc: func(ctx context.Context, id int) (*models.Booking, error) {
				return &models.Booking{ID: id, Name: "Test User"}, nil
			},
			mockDeleteFunc: func(ctx context.Context, id int) error {
				return fmt.Errorf("database error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedErr:    "Failed to delete booking",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &MockBookingRepository{
				GetByIDFunc: tc.mockGetByIDFunc,
				DeleteFunc:  tc.mockDeleteFunc,
			}

			// Create handler with mock
			handler := handlers.NewBookingHandler(mockRepo)

			// Create request with URL parameter
			req := httptest.NewRequest("DELETE", "/api/v1/bookings/"+tc.bookingID, nil)

			// Setup chi context with URL parameters
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.bookingID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Create response recorder
			w := httptest.NewRecorder()

			// Call handler
			handler.Delete(w, req)

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

func TestResponseHeaders(t *testing.T) {

	log.Println("Starting TestTestResponseHeaders")
	// Create mock repository
	mockRepo := &MockBookingRepository{
		GetAllFunc: func(ctx context.Context, includeArchived bool) ([]*models.Booking, error) {
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

func TestArchiveBookingHandler(t *testing.T) {
	tests := []struct {
		name            string
		bookingID       string
		mockGetByIDFunc func(context.Context, int) (*models.Booking, error)
		mockArchiveFunc func(context.Context, int) error
		expectedStatus  int
		expectedErr     string
	}{
		{
			name:      "Successfully archive booking",
			bookingID: "123",
			mockGetByIDFunc: func(ctx context.Context, id int) (*models.Booking, error) {
				return &models.Booking{ID: id, Name: "Test User", Archived: false}, nil
			},
			mockArchiveFunc: func(ctx context.Context, id int) error {
				return nil
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:      "Invalid booking ID",
			bookingID: "abc",
			mockGetByIDFunc: func(ctx context.Context, id int) (*models.Booking, error) {
				return nil, nil // Should not be called
			},
			mockArchiveFunc: func(ctx context.Context, id int) error {
				return nil // Should not be called
			},
			expectedStatus: http.StatusBadRequest,
			expectedErr:    "Invalid booking ID",
		},
		{
			name:      "Booking not found",
			bookingID: "456",
			mockGetByIDFunc: func(ctx context.Context, id int) (*models.Booking, error) {
				return nil, nil
			},
			mockArchiveFunc: func(ctx context.Context, id int) error {
				return nil // Should not be called
			},
			expectedStatus: http.StatusNotFound,
			expectedErr:    "Booking not found",
		},
		{
			name:      "Already archived",
			bookingID: "789",
			mockGetByIDFunc: func(ctx context.Context, id int) (*models.Booking, error) {
				return &models.Booking{ID: id, Name: "Test User", Archived: true}, nil
			},
			mockArchiveFunc: func(ctx context.Context, id int) error {
				return nil
			},
			expectedStatus: http.StatusNoContent, // Idempotent operation
		},
		{
			name:      "Database error",
			bookingID: "101",
			mockGetByIDFunc: func(ctx context.Context, id int) (*models.Booking, error) {
				return &models.Booking{ID: id, Name: "Test User", Archived: false}, nil
			},
			mockArchiveFunc: func(ctx context.Context, id int) error {
				return fmt.Errorf("database error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedErr:    "Failed to archive booking",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &MockBookingRepository{
				GetByIDFunc: tc.mockGetByIDFunc,
				ArchiveFunc: tc.mockArchiveFunc,
			}

			// Create handler with mock
			handler := handlers.NewBookingHandler(mockRepo)

			// Create request with URL parameter
			req := httptest.NewRequest("POST", "/api/v1/bookings/"+tc.bookingID+"/archive", nil)

			// Setup chi context with URL parameters
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.bookingID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Create response recorder
			w := httptest.NewRecorder()

			// Call handler
			handler.Archive(w, req)

			// Check status code
			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			// For valid ID, verify Archive was called
			if tc.expectedStatus == http.StatusNoContent {
				id, _ := strconv.Atoi(tc.bookingID)

				// Only check if Archive was called for non-archived bookings
				if tc.name != "Already archived" {
					if !mockRepo.ArchiveCalled {
						t.Errorf("Expected Archive to be called for ID %d, but it wasn't", id)
					}
					if mockRepo.ArchiveArg != id {
						t.Errorf("Archive called with wrong ID, expected %d, got %d", id, mockRepo.ArchiveArg)
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

func TestUnarchiveBookingHandler(t *testing.T) {
	log.Println("Starting TestUnarchiveBookingHandler")
	tests := []struct {
		name              string
		bookingID         string
		mockGetByIDFunc   func(context.Context, int) (*models.Booking, error)
		mockUnarchiveFunc func(context.Context, int) error
		expectedStatus    int
		expectedErr       string
	}{
		{
			name:      "Successfully unarchive booking",
			bookingID: "123",
			mockGetByIDFunc: func(ctx context.Context, id int) (*models.Booking, error) {
				return &models.Booking{ID: id, Name: "Test User", Archived: true}, nil
			},
			mockUnarchiveFunc: func(ctx context.Context, id int) error {
				return nil
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:      "Invalid booking ID",
			bookingID: "abc",
			mockGetByIDFunc: func(ctx context.Context, id int) (*models.Booking, error) {
				return nil, nil // Should not be called
			},
			mockUnarchiveFunc: func(ctx context.Context, id int) error {
				return nil // Should not be called
			},
			expectedStatus: http.StatusBadRequest,
			expectedErr:    "Invalid booking ID",
		},
		{
			name:      "Booking not found",
			bookingID: "456",
			mockGetByIDFunc: func(ctx context.Context, id int) (*models.Booking, error) {
				return nil, nil
			},
			mockUnarchiveFunc: func(ctx context.Context, id int) error {
				return nil // Should not be called
			},
			expectedStatus: http.StatusNotFound,
			expectedErr:    "Booking not found",
		},
		{
			name:      "Already active (not archived)",
			bookingID: "789",
			mockGetByIDFunc: func(ctx context.Context, id int) (*models.Booking, error) {
				return &models.Booking{ID: id, Name: "Test User", Archived: false}, nil
			},
			mockUnarchiveFunc: func(ctx context.Context, id int) error {
				return nil
			},
			expectedStatus: http.StatusNoContent, // Idempotent operation
		},
		{
			name:      "Database error",
			bookingID: "101",
			mockGetByIDFunc: func(ctx context.Context, id int) (*models.Booking, error) {
				return &models.Booking{ID: id, Name: "Test User", Archived: true}, nil
			},
			mockUnarchiveFunc: func(ctx context.Context, id int) error {
				return fmt.Errorf("database error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedErr:    "Failed to unarchive booking",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &MockBookingRepository{
				GetByIDFunc:   tc.mockGetByIDFunc,
				UnarchiveFunc: tc.mockUnarchiveFunc,
			}

			// Create handler with mock
			handler := handlers.NewBookingHandler(mockRepo)

			// Create request with URL parameter
			req := httptest.NewRequest("POST", "/api/v1/bookings/"+tc.bookingID+"/unarchive", nil)

			// Setup chi context with URL parameters
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.bookingID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Create response recorder
			w := httptest.NewRecorder()

			// Call handler
			handler.Unarchive(w, req)

			// Check status code
			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			// For valid ID, verify Unarchive was called
			if tc.expectedStatus == http.StatusNoContent {
				id, _ := strconv.Atoi(tc.bookingID)

				// Only check if Unarchive was called for archived bookings
				if tc.name != "Already active (not archived)" {
					if !mockRepo.UnarchiveCalled {
						t.Errorf("Expected Unarchive to be called for ID %d, but it wasn't", id)
					}
					if mockRepo.UnarchiveArg != id {
						t.Errorf("Unarchive called with wrong ID, expected %d, got %d", id, mockRepo.UnarchiveArg)
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

func TestGetAllBookingsWithArchiveFiltering(t *testing.T) {
	tests := []struct {
		name             string
		queryParams      string
		mockGetAllFunc   func(context.Context, bool) ([]*models.Booking, error)
		expectedStatus   int
		expectedCount    int
		expectedArchived bool
	}{
		{
			name:        "Get active bookings only (default)",
			queryParams: "",
			mockGetAllFunc: func(ctx context.Context, includeArchived bool) ([]*models.Booking, error) {
				if includeArchived {
					t.Error("Expected includeArchived=false, got true")
				}
				return []*models.Booking{
					{ID: 1, Name: "Active 1", Archived: false},
					{ID: 2, Name: "Active 2", Archived: false},
				}, nil
			},
			expectedStatus:   http.StatusOK,
			expectedCount:    2,
			expectedArchived: false,
		},
		{
			name:        "Get active bookings only (explicit)",
			queryParams: "?include_archived=false",
			mockGetAllFunc: func(ctx context.Context, includeArchived bool) ([]*models.Booking, error) {
				if includeArchived {
					t.Error("Expected includeArchived=false, got true")
				}
				return []*models.Booking{
					{ID: 1, Name: "Active 1", Archived: false},
					{ID: 2, Name: "Active 2", Archived: false},
				}, nil
			},
			expectedStatus:   http.StatusOK,
			expectedCount:    2,
			expectedArchived: false,
		},
		{
			name:        "Get all bookings including archived",
			queryParams: "?include_archived=true",
			mockGetAllFunc: func(ctx context.Context, includeArchived bool) ([]*models.Booking, error) {
				if !includeArchived {
					t.Error("Expected includeArchived=true, got false")
				}
				return []*models.Booking{
					{ID: 1, Name: "Active 1", Archived: false},
					{ID: 2, Name: "Active 2", Archived: false},
					{ID: 3, Name: "Archived 1", Archived: true},
					{ID: 4, Name: "Archived 2", Archived: true},
				}, nil
			},
			expectedStatus:   http.StatusOK,
			expectedCount:    4,
			expectedArchived: true,
		},
		{
			name:        "Invalid include_archived parameter",
			queryParams: "?include_archived=invalid",
			mockGetAllFunc: func(ctx context.Context, includeArchived bool) ([]*models.Booking, error) {
				// Should default to false for invalid values
				if includeArchived {
					t.Error("Expected includeArchived=false for invalid parameter, got true")
				}
				return []*models.Booking{
					{ID: 1, Name: "Active 1", Archived: false},
					{ID: 2, Name: "Active 2", Archived: false},
				}, nil
			},
			expectedStatus:   http.StatusOK,
			expectedCount:    2,
			expectedArchived: false,
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

			// Create request with query parameters
			req := httptest.NewRequest("GET", "/api/v1/bookings"+tc.queryParams, nil)
			w := httptest.NewRecorder()

			// Call handler
			handler.GetAll(w, req)

			// Check status code
			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			// Verify GetAll was called with correct includeArchived parameter
			if mockRepo.GetAllCalled && mockRepo.GetAllIncludeArchived != tc.expectedArchived {
				t.Errorf("Expected GetAll called with includeArchived=%v, got %v",
					tc.expectedArchived, mockRepo.GetAllIncludeArchived)
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
		})
	}
}
