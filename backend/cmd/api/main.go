package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/config"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/database"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/handlers"
	custommiddleware "github.com/joshuagudgel/toasted-coffee/backend/internal/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// run database migrations
	log.Println("Running database migrations...")
	migrationSQL1, err := os.ReadFile("internal/database/migrations/01_create_bookings_table.sql")
	if err != nil {
		log.Printf("Warning: Could not read migration file: %v", err)
	} else {
		_, err = db.Pool.Exec(context.Background(), string(migrationSQL1))
		if err != nil {
			// Check if error is because table already exists (which is fine)
			if strings.Contains(err.Error(), "already exists") {
				log.Println("Tables already exist, skipping migration")
			} else {
				log.Printf("Warning: Migration 1 error: %v", err)
			}
		} else {
			log.Println("Migration 1 executed successfully")
		}
	}

	migrationSQL2, err := os.ReadFile("internal/database/migrations/02_create_users_table.sql")
	if err == nil {
		_, err := db.Pool.Exec(context.Background(), string(migrationSQL2))
		if err != nil {
			// Check if error is because table already exists (which is fine)
			if strings.Contains(err.Error(), "already exists") {
				log.Println("Tables already exist, skipping migration")
			} else {
				log.Printf("Warning: Migration 2 error: %v", err)
			}
		} else {
			log.Println("Migration 2 executed successfully")
		}
	} else {
		log.Printf("Warning: Could not read migration file: %v", err)
	}

	log.Println("Setting up admin user...")
	passwordHash := os.Getenv("ADMIN_USER_PASSWORD_HASH")
	if passwordHash == "" {
		log.Println("Warning: ADMIN_USER_PASSWORD_HASH not set")
	}

	// Use parameterized query for inserting admin user
	_, err = db.Pool.Exec(context.Background(), `
    	INSERT INTO users (username, password, role) 
    	VALUES ($1, $2, $3) 
    	ON CONFLICT (username) DO UPDATE SET 
    	    password = $2,
    	    role = $3
`, "admin", passwordHash, "admin")
	if err != nil {
		log.Printf("Warning: Failed to create admin user: %v", err)
	} else {
		log.Println("Admin user created successfully")
	}

	// Initialize repositories
	bookingRepo := database.NewBookingRepository(db)
	userRepo := database.NewUserRepository(db)

	// Initialize handlers
	bookingHandler := handlers.NewBookingHandler(bookingRepo)
	authHandler := handlers.NewAuthHandler(userRepo)

	// Initialize router
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(custommiddleware.CORS(cfg.AllowOrigins))

	// Routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes (no auth required)
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/refresh", authHandler.RefreshToken)
		r.Post("/bookings", bookingHandler.Create)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(custommiddleware.JWTAuth)

			// Bookings
			r.Get("/bookings", bookingHandler.GetAll)
			r.Get("/bookings/{id}", bookingHandler.GetByID)

			// Auth validation
			r.Get("/auth/validate", authHandler.ValidateToken)
		})
	})

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Check Database Connection
	r.Get("/api/v1/db-check", func(w http.ResponseWriter, r *http.Request) {
		// Test basic connectivity
		if err := db.Pool.Ping(r.Context()); err != nil {
			log.Printf("Database ping failed: %v", err)
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}

		// Test a simple query
		var result int
		err := db.Pool.QueryRow(r.Context(), "SELECT 1").Scan(&result)
		if err != nil || result != 1 {
			log.Printf("Database query test failed: %v", err)
			http.Error(w, "Database query error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Database connection OK"))
	})

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
