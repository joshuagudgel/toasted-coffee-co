package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/database"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/models"
)

// BookingHandler handles HTTP requests related to bookings
type BookingHandler struct {
	repo *database.BookingRepository
}

// NewBookingHandler creates a new booking handler
func NewBookingHandler(repo *database.BookingRepository) *BookingHandler {
	return &BookingHandler{repo: repo}
}

// Create handles creation of a new booking
func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var booking models.Booking
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.repo.Create(r.Context(), &booking)
	if err != nil {
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
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	booking, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to retrieve booking", http.StatusInternalServerError)
		return
	}

	if booking == nil {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(booking)
}

// GetAll retrieves all bookings
func (h *BookingHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching all bookings")

	bookings, err := h.repo.GetAll(r.Context())
	if err != nil {
		log.Printf("ERROR in GetAll: %v", err) // Log the actual error
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
