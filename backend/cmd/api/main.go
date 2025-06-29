package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate" // Add this import
	"github.com/joshuagudgel/toasted-coffee/backend/internal/config"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/database"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/handlers"
	custommiddleware "github.com/joshuagudgel/toasted-coffee/backend/internal/middleware"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/services"
	"golang.org/x/crypto/bcrypt"
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

	// Log database migrations
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

	migrationSQL3, err := os.ReadFile("internal/database/migrations/03_add_contact_fields.sql")
	if err == nil {
		_, err := db.Pool.Exec(context.Background(), string(migrationSQL3))
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				log.Println("Columns already exist, skipping migration")
			} else {
				log.Printf("Warning: Migration 3 error: %v", err)
			}
		} else {
			log.Println("Migration 3 executed successfully")
		}
	} else {
		log.Printf("Warning: Could not read migration file: %v", err)
	}

	migrationSQL4, err := os.ReadFile("internal/database/migrations/04_add_archived_column.sql")
	if err == nil {
		_, err := db.Pool.Exec(context.Background(), string(migrationSQL4))
		if err != nil {
			log.Printf("Warning: Migration 4 error: %v", err)
		} else {
			log.Println("Migration 4 executed successfully")
		}
	} else {
		log.Printf("Warning: Could not read migration file: %v", err)
	}

	migrationSQL5, err := os.ReadFile("internal/database/migrations/05_create_menu_items_table.sql")
	if err == nil {
		_, err := db.Pool.Exec(context.Background(), string(migrationSQL5))
		if err != nil {
			log.Printf("Warning: Migration 5 error: %v", err)
		} else {
			log.Println("Migration 5 executed successfully")
		}
	} else {
		log.Printf("Warning: Could not read migration file: %v", err)
	}

	migrationSQL6, err := os.ReadFile("internal/database/migrations/06_add_menu_items.sql")
	if err == nil {
		_, err := db.Pool.Exec(context.Background(), string(migrationSQL6))
		if err != nil {
			log.Printf("Warning: Migration 6 error: %v", err)
		} else {
			log.Println("Migration 6 executed successfully")
		}
	} else {
		log.Printf("Warning: Could not read migration file: %v", err)
	}

	log.Println("Setting up admin user...")
	var count int
	err = db.Pool.QueryRow(context.Background(), `
	    SELECT COUNT(*) FROM users WHERE username = $1
	`, "admin").Scan(&count)

	// Replace your existing admin user creation code with this:
	if err != nil {
		log.Printf("Warning: Failed to check for admin user: %v", err)
	} else if count == 0 {
		// Generate a fresh bcrypt hash directly in the code
		rawHash, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Warning: Failed to generate hash: %v", err)
		} else {
			// Insert with the freshly generated hash
			_, err = db.Pool.Exec(context.Background(), `
            INSERT INTO users (username, password, role) 
            VALUES ($1, $2, $3)
        `, "admin", string(rawHash), "admin")

			if err != nil {
				log.Printf("Warning: Failed to create admin user: %v", err)
			} else {
				// Verify what was stored
				var storedHash string
				err = db.Pool.QueryRow(context.Background(), `
                SELECT password FROM users WHERE username = $1
            `, "admin").Scan(&storedHash)
				if err != nil {
					log.Printf("Warning: Failed to retrieve hash: %v", err)
				} else {
					log.Printf("Admin user created successfully")
					log.Printf("DEBUG - Generated hash: %s", string(rawHash))
					log.Printf("DEBUG - Stored hash: %s", storedHash)
				}
			}
		}
	} else {
		log.Println("Admin user already exists, skipping password update")
	}

	// Initialize repositories
	bookingRepo := database.NewBookingRepository(db)
	userRepo := database.NewUserRepository(db)
	menuRepo := database.NewMenuRepository(db)

	// First, create the email service as a shared resource
	emailService := services.NewEmailService()

	// Initialize handlers
	bookingHandler := handlers.NewBookingHandler(bookingRepo, emailService)
	authHandler := handlers.NewAuthHandler(userRepo)
	menuHandler := handlers.NewMenuHandler(menuRepo)

	// Create the contact handler with the email service
	contactHandler := handlers.NewContactHandler(emailService)

	// Initialize router
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(custommiddleware.CORS(cfg.AllowOrigins))

	// Routes
	r.Route("/api/v1", func(r chi.Router) {
		// Auth endpoints
		r.Group(func(r chi.Router) {
			r.Use(httprate.LimitByIP(20, 1*time.Minute))
			r.Post("/auth/login", authHandler.Login)
			r.Post("/auth/refresh", authHandler.RefreshToken)
			r.Post("/auth/logout", authHandler.Logout)
		})

		// Booking creation - limit to 10 requests per minute
		// This prevents spam bookings while allowing legitimate use
		r.With(httprate.LimitByIP(10, 1*time.Minute)).Post("/bookings", bookingHandler.Create)

		// Contact/inquiry endpoint - limit to 5 requests per minute
		// This is more restricted since it's a common target for spam
		r.With(httprate.LimitByIP(5, 1*time.Minute)).Post("/contact", contactHandler.HandleInquiry)

		// Public read-only endpoints - higher limits (30 per minute)
		// These are less sensitive and used more frequently by legitimate users
		r.Group(func(r chi.Router) {
			r.Use(httprate.LimitByIP(30, 1*time.Minute))
			r.Get("/menu", menuHandler.GetAll)
			r.Get("/menu/{type}", menuHandler.GetByType)
		})

		// Protected routes - admin functions with JWT auth
		r.Group(func(r chi.Router) {
			r.Use(custommiddleware.JWTAuth)

			// You can also rate limit protected routes, but with higher thresholds
			// This prevents API abuse even from authenticated users
			r.Use(httprate.LimitByIP(60, 1*time.Minute))

			// Bookings management
			r.Get("/bookings", bookingHandler.GetAll)
			r.Get("/bookings/{id}", bookingHandler.GetByID)
			r.Delete("/bookings/{id}", bookingHandler.Delete)
			r.Put("/bookings/{id}", bookingHandler.Update)
			r.Post("/bookings/{id}/archive", bookingHandler.Archive)
			r.Post("/bookings/{id}/unarchive", bookingHandler.Unarchive)

			// Menu management
			r.Post("/menu", menuHandler.Create)
			r.Put("/menu/{id}", menuHandler.Update)
			r.Delete("/menu/{id}", menuHandler.Delete)

			// Auth validation
			r.Get("/auth/validate", authHandler.ValidateToken)
		})
	})

	// Health check with minimal rate limiting
	r.With(httprate.LimitByIP(10, 1*time.Minute)).Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// DB check with minimal rate limiting
	r.With(httprate.LimitByIP(10, 1*time.Minute)).Get("/api/v1/db-check", func(w http.ResponseWriter, r *http.Request) {
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
