package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/joshuagudgel/toasted-coffee/backend/internal/auth"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/database"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo *database.UserRepository
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refreshToken"`
	User         models.User `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshResponse struct {
	AccessToken string `json:"accessToken"`
}

func NewAuthHandler(userRepo *database.UserRepository) *AuthHandler {
	return &AuthHandler{userRepo: userRepo}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("Login request received for user: %s", req.Username)

	// Get user by username
	user, err := h.userRepo.GetByUsername(r.Context(), req.Username)
	if err != nil {
		log.Printf("ERROR: User '%s' lookup failed: %v", req.Username, err)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	log.Printf("User found: %s (ID: %d, Role: %s)", user.Username, user.ID, user.Role)

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Printf("ERROR: Password verification failed for '%s': %v", user.Username, err)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	log.Printf("Password verification successful for user: %s", user.Username)

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	log.Printf("JWT token generated successfully")

	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		log.Printf("ERROR: Refresh token generation failed: %v", err)
		http.Error(w, "Error generating refresh token", http.StatusInternalServerError)
		return
	}
	log.Printf("Refresh token generated successfully")

	// Return token and user info
	user.Password = "" // Don't send password back

	w.Header().Set("Content-Type", "application/json")
	resp := LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         *user,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("ERROR: Failed to encode login response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	log.Printf("Login successful for user: %s, role: %s", user.Username, user.Role)
}

func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	// The JWT middleware already validated the token
	// Just extract the claims and return user data

	// In a real implementation, you would extract claims from context
	claims, _ := auth.ExtractClaimsFromContext(r.Context())

	// Return user info based on claims
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"userId": claims.UserID,
		"role":   claims.Role,
	})
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate refresh token
	userID, err := auth.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Get user details to include role information
	user, err := h.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Generate new access token
	newAccessToken, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RefreshResponse{
		AccessToken: newAccessToken,
	})
}
