package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/database"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/models"
)

// BookingHandler handles HTTP requests related to bookings
type BookingHandler struct {
	repo database.BookingRepositoryInterface // Changed from *database.BookingRepository
}

// NewBookingHandler creates a new booking handler
func NewBookingHandler(repo database.BookingRepositoryInterface) *BookingHandler {
	return &BookingHandler{repo: repo}
}

// Create handles creation of a new booking
func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var booking models.Booking

	// Log the incoming request
	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	log.Printf("Received booking request: %s", string(body))

	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the booking email and phone
	if booking.Email == "" && booking.Phone == "" {
		log.Println("Booking rejected: no contact information provided")
		http.Error(w, "Email or phone number is required", http.StatusBadRequest)
		return
	}

	_, err := time.Parse("2006-01-02", booking.Date)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	// Log the decoded booking
	log.Printf("Decoded booking: %+v", booking)

	id, err := h.repo.Create(r.Context(), &booking)
	if err != nil {
		log.Printf("Error creating booking: %v", err)
		http.Error(w, "Failed to create booking", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"message": "Booking created successfully",
	})
}

// GetByID retrieves a booking by ID
func (h *BookingHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Parse the ID from the URL
	idStr := chi.URLParam(r, "id")
	log.Printf("GetByID request for booking: %s", idStr)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Handle invalid ID format specifically
		log.Printf("Invalid booking ID format: %s", idStr)
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	// Get the booking from the repository
	booking, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("Error retrieving booking %d: %v", id, err)

		// Check for "not found" error specifically
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Booking not found", http.StatusNotFound)
			return
		}

		// Return 500 for other errors
		http.Error(w, "Failed to retrieve booking", http.StatusInternalServerError)
		return
	}

	// Check if booking is nil even without an error
	if booking == nil {
		log.Printf("Booking not found with ID: %d", id)
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	// Return the booking as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(booking); err != nil {
		log.Printf("Error encoding booking response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetAll retrieves all bookings
func (h *BookingHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	includeArchived := r.URL.Query().Get("include_archived") == "true"
	log.Printf("Fetching bookings, includeArchived: %v", includeArchived)

	bookings, err := h.repo.GetAll(r.Context(), includeArchived)
	if err != nil {
		log.Printf("ERROR in GetAll: %v", err)
		http.Error(w, "Failed to retrieve bookings", http.StatusInternalServerError)
		return
	}

	log.Printf("Found %d bookings", len(bookings))

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(bookings); err != nil {
		log.Printf("ERROR encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Delete removes a booking
func (h *BookingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Parse booking ID from the URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid booking ID format: %s", idStr)
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	// Check if the booking exists first
	booking, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("Error checking booking existence %d: %v", id, err)

		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Booking not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to check booking", http.StatusInternalServerError)
		return
	}

	// If booking is nil, it doesn't exist
	if booking == nil {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	// Delete the booking
	err = h.repo.Delete(r.Context(), id)
	if err != nil {
		log.Printf("Error deleting booking %d: %v", id, err)

		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Booking not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to delete booking", http.StatusInternalServerError)
		return
	}

	// Return success with no content
	w.WriteHeader(http.StatusNoContent) // 204 status code indicates successful deletion with no content to return
}

// Update modifies an existing booking
func (h *BookingHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Parse booking ID from the URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid booking ID format: %s", idStr)
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	// Get current booking to check for archive status changes
	currentBooking, err := h.repo.GetByID(r.Context(), id)
	if err != nil || currentBooking == nil {
		log.Printf("Cannot find booking to update: %d", id)
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	// Parse request body
	var booking models.Booking

	// Log the incoming request
	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	log.Printf("Received booking update request for ID %d: %s", id, string(body))

	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate booking data (same validation as Create)
	if booking.Email == "" && booking.Phone == "" {
		log.Println("Booking update rejected: no contact information provided")
		http.Error(w, "Email or phone number is required", http.StatusBadRequest)
		return
	}

	_, err = time.Parse("2006-01-02", booking.Date)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	if len(booking.CoffeeFlavors) < 1 {
		http.Error(w, "At least one coffee flavor is required", http.StatusBadRequest)
		return
	}

	if len(booking.MilkOptions) < 1 {
		http.Error(w, "At least one milk option is required", http.StatusBadRequest)
		return
	}

	// Track archive status changes
	if currentBooking.Archived != booking.Archived {
		if booking.Archived {
			log.Printf("Booking %d is being archived via update", id)
		} else {
			log.Printf("Booking %d is being unarchived via update", id)
		}
	}

	// Update the booking
	err = h.repo.Update(r.Context(), id, &booking)
	if err != nil {
		log.Printf("Error updating booking %d: %v", id, err)

		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Booking not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to update booking", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully updated booking %d (archived status: %v)", id, booking.Archived)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Booking updated successfully",
	})
}

// Archive marks a booking as archived
func (h *BookingHandler) Archive(w http.ResponseWriter, r *http.Request) {
	// Parse booking ID from the URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid booking ID format for archive: %s", idStr)
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	// Check if booking exists first
	booking, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("Error checking booking existence %d: %v", id, err)
		http.Error(w, "Failed to check booking", http.StatusInternalServerError)
		return
	}

	if booking == nil {
		log.Printf("Cannot archive non-existent booking: %d", id)
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	// Don't archive if already archived
	if booking.Archived {
		log.Printf("Booking %d is already archived", id)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = h.repo.Archive(r.Context(), id)
	if err != nil {
		log.Printf("Error archiving booking %d: %v", id, err)

		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Booking not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to archive booking", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully archived booking %d", id)
	w.WriteHeader(http.StatusNoContent)
}

// Unarchive marks a booking as unarchived
func (h *BookingHandler) Unarchive(w http.ResponseWriter, r *http.Request) {
	// Parse booking ID from the URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid booking ID format for unarchive: %s", idStr)
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	// Check if booking exists first
	booking, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("Error checking booking existence %d: %v", id, err)
		http.Error(w, "Failed to check booking", http.StatusInternalServerError)
		return
	}

	if booking == nil {
		log.Printf("Cannot unarchive non-existent booking: %d", id)
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	// Don't unarchive if already active
	if !booking.Archived {
		log.Printf("Booking %d is already active (not archived)", id)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = h.repo.Unarchive(r.Context(), id)
	if err != nil {
		log.Printf("Error unarchiving booking %d: %v", id, err)

		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Booking not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to unarchive booking", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully unarchived booking %d", id)
	w.WriteHeader(http.StatusNoContent)
}
