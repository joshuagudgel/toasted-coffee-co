package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// Make this exportable so middleware can use it
type ContextKey string

// Export this key for middleware to use
const ClaimsContextKey ContextKey = "claims"

// Predefined errors for more secure error handling
var (
	ErrTokenExpired     = errors.New("token has expired")
	ErrTokenNotValidYet = errors.New("token not valid yet")
	ErrTokenInvalid     = errors.New("token is invalid")
	ErrMissingSecret    = errors.New("jwt secret key is not configured")
)

// Token-related functions and structures
type Claims struct {
	UserID int    `json:"userId"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Global variables
var secretKey []byte
var refreshSecretKey []byte
var tokenExpiry time.Duration
var refreshTokenExpiry time.Duration

func init() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
	// Load the secret key from environment
	secretKeyStr := os.Getenv("JWT_SECRET")
	if secretKeyStr == "" {
		// TODO: remove
		log.Println("WARNING: JWT_SECRET environment variable not set! Using a random key for this session.")

		// Generate a random key for this session
		randomKey := uuid.New().String()
		secretKey = []byte(randomKey)
	} else {
		secretKey = []byte(secretKeyStr)
		log.Printf("JWT secret key configured with length: %d bytes", len(secretKey))
	}

	// Add this for refreshSecretKey
	refreshSecretKeyStr := os.Getenv("JWT_REFRESH_SECRET")
	if refreshSecretKeyStr == "" {
		log.Println("WARNING: JWT_REFRESH_SECRET environment variable not set! Using same key as JWT_SECRET.")
		refreshSecretKey = secretKey
	} else {
		refreshSecretKey = []byte(refreshSecretKeyStr)
		log.Printf("JWT refresh secret key configured with length: %d bytes", len(refreshSecretKey))
	}

	// Parse token expiry from environment
	tokenExpiryStr := os.Getenv("TOKEN_EXPIRY")
	if tokenExpiryStr == "" {
		log.Println("WARNING: TOKEN_EXPIRY not set, defaulting to 15m")
		tokenExpiry = 15 * time.Minute
	} else {
		parsed, err := time.ParseDuration(tokenExpiryStr)
		if err != nil {
			log.Printf("WARNING: Invalid TOKEN_EXPIRY format: %v, defaulting to 15m", err)
			tokenExpiry = 15 * time.Minute
		} else {
			tokenExpiry = parsed
		}
	}
	log.Printf("Access token expiry set to: %s", tokenExpiry)

	// Parse refresh token expiry from environment
	refreshTokenExpiryStr := os.Getenv("REFRESH_TOKEN_EXPIRY")
	if refreshTokenExpiryStr == "" {
		log.Println("WARNING: REFRESH_TOKEN_EXPIRY not set, defaulting to 7d")
		refreshTokenExpiry = 7 * 24 * time.Hour
	} else {
		parsed, err := time.ParseDuration(refreshTokenExpiryStr)
		if err != nil {
			log.Printf("WARNING: Invalid REFRESH_TOKEN_EXPIRY format: %v, defaulting to 7d", err)
			refreshTokenExpiry = 7 * 24 * time.Hour
		} else {
			refreshTokenExpiry = parsed
		}
	}
	log.Printf("Refresh token expiry set to: %s", refreshTokenExpiry)
}

// Token generation and validation functions
func GenerateToken(userID int, role string) (string, error) {
	// Create unique token ID
	tokenID := uuid.New().String()

	// Define accepted audiences
	audiences := []string{"toasted-coffee-admin", "toasted-coffee-api"}

	// Create claims with expiration time and additional security claims
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "toasted-coffee-co",
			Audience:  audiences,
			ID:        tokenID,
		},
	}

	// Generate token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	return token.SignedString(secretKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		// Convert JWT errors to our custom errors for more secure error handling
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotValidYet
		}

		// Don't expose specific JWT errors to callers
		log.Printf("JWT validation error (not exposed to client): %v", err)
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}

	// Explicitly check expiration even though the JWT library does this
	// This is for clarity and additional security
	now := time.Now()
	if now.After(claims.ExpiresAt.Time) {
		return nil, ErrTokenExpired
	}

	// Explicitly check not-before time
	if now.Before(claims.NotBefore.Time) {
		return nil, ErrTokenNotValidYet
	}

	// Verify audience - token must be intended for our service
	validAudience := false
	for _, audience := range claims.Audience {
		if audience == "toasted-coffee-api" || audience == "toasted-coffee-admin" {
			validAudience = true
			break
		}
	}

	if !validAudience {
		return nil, errors.New("token has invalid audience")
	}

	return claims, nil
}

func ExtractClaimsFromContext(ctx context.Context) (*Claims, bool) {
	if claims, ok := ctx.Value(ClaimsContextKey).(*Claims); ok {
		return claims, true
	}
	return nil, false
}

// Refresh token functionality
func GenerateRefreshToken(userID int) (string, error) {
	// Create unique token ID for revocation capability
	tokenID := uuid.New().String()

	refreshClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenExpiry)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "toasted-coffee-co",
		Subject:   fmt.Sprintf("%d", userID),
		Audience:  []string{"toasted-coffee-refresh"},
		ID:        tokenID,
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	return refreshToken.SignedString(refreshSecretKey)
}

func ValidateRefreshToken(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return refreshSecretKey, nil
	})

	if err != nil {
		// Convert JWT errors to our custom errors
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, ErrTokenExpired
		}
		// Don't expose specific JWT errors
		log.Printf("Refresh token validation error (not exposed): %v", err)
		return 0, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return 0, ErrTokenInvalid
	}

	// Explicitly check expiration
	now := time.Now()
	if now.After(claims.ExpiresAt.Time) {
		return 0, ErrTokenExpired
	}

	// Explicitly check not-before time
	if now.Before(claims.NotBefore.Time) {
		return 0, ErrTokenNotValidYet
	}

	// Verify this is a refresh token
	validAudience := false
	for _, audience := range claims.Audience {
		if audience == "toasted-coffee-refresh" {
			validAudience = true
			break
		}
	}

	if !validAudience {
		return 0, errors.New("token has invalid audience")
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, errors.New("invalid user ID in token")
	}

	return userID, nil
}

// IsAdmin helper function to check if a user has admin role
func IsAdmin(claims *Claims) bool {
	return claims != nil && claims.Role == "admin"
}
