package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

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
	startTime := time.Now()
	log.Printf("LOGIN START: Authentication request received at %v", startTime)

	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("LOGIN TIMING: Request body decoded in %v", time.Since(startTime))

	log.Printf("Login request received for user: %s", req.Username)

	// Get user by username
	userLookupStart := time.Now()
	user, err := h.userRepo.GetByUsername(r.Context(), req.Username)
	if err != nil {
		log.Printf("ERROR: User '%s' lookup failed: %v", req.Username, err)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	log.Printf("LOGIN TIMING: Database user lookup took %v", time.Since(userLookupStart))
	log.Printf("User found: %s (ID: %d, Role: %s)", user.Username, user.ID, user.Role)

	// Compare passwords
	pwCompareStart := time.Now()
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Printf("ERROR: Password verification failed for '%s': %v", user.Username, err)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	log.Printf("LOGIN TIMING: Password verification took %v", time.Since(pwCompareStart))
	log.Printf("Password verification successful for user: %s", user.Username)

	// Generate JWT token
	tokenGenStart := time.Now()
	token, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	log.Printf("LOGIN TIMING: JWT token generation took %v", time.Since(tokenGenStart))
	log.Printf("JWT token generated successfully")

	refreshTokenStart := time.Now()
	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		log.Printf("ERROR: Refresh token generation failed: %v", err)
		http.Error(w, "Error generating refresh token", http.StatusInternalServerError)
		return
	}
	log.Printf("LOGIN TIMING: Refresh token generation took %v", time.Since(refreshTokenStart))
	log.Printf("Refresh token generated successfully")

	// Set secure HttpOnly cookies instead of returning tokens in response body
	// Access token cookie - shorter expiration
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil, // true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400, // 24 hours
	})

	// Refresh token cookie - longer expiration
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/api/v1/auth/refresh", // Restrict to refresh endpoint only
		HttpOnly: true,
		Secure:   r.TLS != nil, // true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
		MaxAge:   259200, // 3 days
	})

	// Return user info only (without tokens)
	w.Header().Set("Content-Type", "application/json")
	resp := struct {
		User models.User `json:"user"`
	}{
		User: *user,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("ERROR: Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("LOGIN COMPLETE: Total authentication time: %v", time.Since(startTime))
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
	// Get refresh token from cookie instead of request body
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Refresh token not found", http.StatusUnauthorized)
		return
	}

	// Validate refresh token
	userID, err := auth.ValidateRefreshToken(refreshCookie.Value)
	if err != nil {
		// Clear the invalid cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/api/v1/auth/refresh",
			HttpOnly: true,
			MaxAge:   -1, // Delete the cookie
		})
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

	// Set new access token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    newAccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400, // 24 hours
	})

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{
		"success": true,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear access token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1, // Delete the cookie
	})

	// Clear refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/v1/auth/refresh",
		HttpOnly: true,
		MaxAge:   -1, // Delete the cookie
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{
		"success": true,
	})
}
