package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/joshuagudgel/toasted-coffee/backend/internal/services"
)

type ContactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

type ContactHandler struct {
	emailService *services.EmailService
}

func NewContactHandler(emailService *services.EmailService) *ContactHandler {
	return &ContactHandler{
		emailService: emailService,
	}
}

func (h *ContactHandler) HandleInquiry(w http.ResponseWriter, r *http.Request) {
	var request ContactRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Error decoding contact request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the request
	if request.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if request.Email == "" && request.Phone == "" {
		http.Error(w, "Email or phone is required", http.StatusBadRequest)
		return
	}

	if request.Message == "" {
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}

	// Send the inquiry email
	err := h.emailService.SendInquiry(
		request.Name,
		request.Email,
		request.Phone,
		request.Message,
	)

	if err != nil {
		log.Printf("Failed to send inquiry email: %v", err)
		http.Error(w, "Failed to send inquiry", http.StatusInternalServerError)
		return
	}

	// Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Inquiry sent successfully",
	})
}
