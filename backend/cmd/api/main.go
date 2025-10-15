package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/config"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/database"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/handlers"
	custommiddleware "github.com/joshuagudgel/toasted-coffee/backend/internal/middleware"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/services"
	"golang.org/x/crypto/bcrypt"
)

// Constants
const (
	// --- Rate Limits ---
	// Public endpoints
	PublicReadLimit  = 100 // requests per minute
	PublicWriteLimit = 10  // requests per minute
	ContactLimit     = 5   // requests per minute (spam protection)

	// Authentication
	AuthLimit = 20 // requests per minute

	// Admin/Protected endpoints
	AdminLimit = 200 // requests per minute (higher for legitimate admin use)

	// Special endpoints
	HealthCheckLimit = 60 // requests per minute (monitoring tools)
)

// Global variables for health check
var (
	serviceStartTime = time.Now()
	lastHealthCheck  time.Time
	healthCheckCount int64
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

	migrationSQL7, err := os.ReadFile("internal/database/migrations/07_create_packages_table.sql")
	if err == nil {
		_, err := db.Pool.Exec(context.Background(), string(migrationSQL7))
		if err != nil {
			log.Printf("Warning: Migration 7 error: %v", err)
		} else {
			log.Println("Migration 7 executed successfully")
		}
	} else {
		log.Printf("Warning: Could not read migration file: %v", err)
	}

	migrationSQL8, err := os.ReadFile("internal/database/migrations/08_add_package_display_order.sql")
	if err == nil {
		_, err := db.Pool.Exec(context.Background(), string(migrationSQL8))
		if err != nil {
			log.Printf("Warning: Migration 8 error: %v", err)
		} else {
			log.Println("Migration 8 executed successfully")
		}
	} else {
		log.Printf("Warning: Could not read migration file: %v", err)
	}

	workingDir, _ := os.Getwd()
	log.Printf("Current working directory: %s", workingDir)
	log.Printf("Attempting to read migration from: %s", filepath.Join(workingDir, "internal/database/migrations/09_add_bookings_outdoor_details.sql"))

	// Check if the file exists before trying to read it
	if _, err := os.Stat(filepath.Join(workingDir, "internal/database/migrations/09_add_bookings_outdoor_details.sql")); os.IsNotExist(err) {
		log.Printf("WARNING: Migration file does not exist at expected path")
	}

	log.Println("Checking database schema information...")
	var tableSchema string
	var tableExists bool
	err = db.Pool.QueryRow(context.Background(), `
    SELECT EXISTS (
        SELECT 1 FROM information_schema.tables 
        WHERE table_name = 'bookings'
    )
`).Scan(&tableExists)
	if err != nil {
		log.Printf("Error checking if bookings table exists: %v", err)
	} else {
		log.Printf("Bookings table exists: %v", tableExists)

		if tableExists {
			err = db.Pool.QueryRow(context.Background(), `
            SELECT table_schema FROM information_schema.tables 
            WHERE table_name = 'bookings' LIMIT 1
        `).Scan(&tableSchema)
			if err != nil {
				log.Printf("Error getting bookings table schema: %v", err)
			} else {
				log.Printf("Bookings table schema: %s", tableSchema)
			}

			// Also check the column structure
			var columns []string
			rows, err := db.Pool.Query(context.Background(), `
            SELECT column_name FROM information_schema.columns 
            WHERE table_name = 'bookings' 
            ORDER BY ordinal_position
        `)
			if err != nil {
				log.Printf("Error fetching column info: %v", err)
			} else {
				defer rows.Close()
				for rows.Next() {
					var col string
					rows.Scan(&col)
					columns = append(columns, col)
				}
				log.Printf("Existing columns: %v", columns)
			}
		}
	}

	// Then load and execute your migration as before
	migrationSQL9, err := os.ReadFile("internal/database/migrations/09_add_bookings_outdoor_details.sql")
	if err == nil {
		// Modify your migration to use the correct schema
		migrationContent := string(migrationSQL9)
		if tableSchema != "" {
			// If we found the schema, add it explicitly to ensure correct targeting
			migrationContent = strings.ReplaceAll(
				migrationContent,
				"ALTER TABLE bookings",
				fmt.Sprintf("ALTER TABLE %s.bookings", tableSchema))
			log.Printf("Updated migration to target schema: %s", tableSchema)
		}

		// Execute the migration
		_, err := db.Pool.Exec(context.Background(), migrationContent)
		if err != nil {
			log.Printf("Warning: Migration 9 error: %v", err)
		} else {
			log.Println("Migration 9 executed successfully")

			// Verify the new columns were added
			var columnsAfter []string
			rows, err := db.Pool.Query(context.Background(), `
            SELECT column_name FROM information_schema.columns 
            WHERE table_name = 'bookings' 
            ORDER BY ordinal_position
        `)
			if err != nil {
				log.Printf("Error fetching post-migration column info: %v", err)
			} else {
				defer rows.Close()
				for rows.Next() {
					var col string
					rows.Scan(&col)
					columnsAfter = append(columnsAfter, col)
				}
				log.Printf("Columns after migration: %v", columnsAfter)

				// Specifically check for our new columns
				hasOutdoor := false
				hasShade := false
				for _, col := range columnsAfter {
					if col == "is_outdoor" {
						hasOutdoor = true
					}
					if col == "has_shade" {
						hasShade = true
					}
				}
				log.Printf("New columns present: is_outdoor=%v, has_shade=%v", hasOutdoor, hasShade)
			}
		}
	} else {
		log.Printf("Warning: Could not read migration file: %v", err)
	}

	migrationSQL10, err := os.ReadFile("internal/database/migrations/10_update_bookings_outdoor_details.sql")
	if err == nil {
		_, err := db.Pool.Exec(context.Background(), string(migrationSQL10))
		if err != nil {
			log.Printf("Warning: Migration 10 error: %v", err)
		} else {
			log.Println("Migration 10 executed successfully")
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
	packageRepo := database.NewPackageRepository(db)

	// First, create the email service as a shared resource
	emailService := services.NewEmailService()

	// Initialize handlers
	bookingHandler := handlers.NewBookingHandler(bookingRepo, emailService)
	authHandler := handlers.NewAuthHandler(userRepo)
	menuHandler := handlers.NewMenuHandler(menuRepo)
	packageHandler := handlers.NewPackageHandler(packageRepo)
	contactHandler := handlers.NewContactHandler(emailService)

	// Initialize router
	mainRouter := chi.NewRouter()

	// Monitoring sub-router without rate limiting
	monitorRouter := chi.NewRouter()
	monitorRouter.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
			"uptime":    time.Since(serviceStartTime).String(),
		})
	})
	monitorRouter.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// API sub-router with rate limiting and authentication
	r := chi.NewRouter()

	// Monitoring endpoints
	// Simple health check endpoint without rate limiting
	// This is useful for uptime monitoring services that need quick checks
	// and should not be blocked by rate limits
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		// Your existing health check code here
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
			"uptime":    time.Since(serviceStartTime).String(),
		})
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Global middleware
	r.Use(custommiddleware.SecureHTTPS)
	r.Use(custommiddleware.SecurityHeaders)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(custommiddleware.CORS(cfg.AllowOrigins))

	// Routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public read-only endpoints
		r.Group(func(r chi.Router) {
			r.Use(httprate.LimitByIP(PublicReadLimit, 1*time.Minute))
			r.Get("/menu", menuHandler.GetAll)
			r.Get("/menu/{type}", menuHandler.GetByType)
			r.Get("/packages", packageHandler.GetAll)
		})

		// Public write endpoints (more restrictive)
		r.Group(func(r chi.Router) {
			r.Use(httprate.LimitByIP(PublicWriteLimit, 1*time.Minute))
			r.Post("/bookings", bookingHandler.Create)
		})

		// Contact endpoint (most restrictive)
		r.With(httprate.LimitByIP(ContactLimit, 1*time.Minute)).
			Post("/contact", contactHandler.HandleInquiry)

		// Authentication endpoints
		r.Group(func(r chi.Router) {
			r.Use(httprate.LimitByIP(AuthLimit, 1*time.Minute))
			r.Post("/auth/login", authHandler.Login)
			r.Post("/auth/refresh", authHandler.RefreshToken)
			r.Post("/auth/logout", authHandler.Logout)
		})

		// Protected admin endpoints
		r.Group(func(r chi.Router) {
			r.Use(custommiddleware.JWTAuth)
			r.Use(httprate.LimitByIP(AdminLimit, 1*time.Minute))

			// All admin endpoints inherit the AdminLimit
			r.Get("/bookings", bookingHandler.GetAll)
			r.Get("/bookings/{id}", bookingHandler.GetByID)
			r.Put("/bookings/{id}", bookingHandler.Update)
			r.Delete("/bookings/{id}", bookingHandler.Delete)
			r.Post("/bookings/{id}/archive", bookingHandler.Archive)
			r.Post("/bookings/{id}/unarchive", bookingHandler.Unarchive)

			r.Post("/menu", menuHandler.Create)
			r.Put("/menu/{id}", menuHandler.Update)
			r.Delete("/menu/{id}", menuHandler.Delete)

			r.Post("/packages", packageHandler.Create)
			r.Get("/packages/{id}", packageHandler.GetByID)
			r.Put("/packages/{id}", packageHandler.Update)
			r.Delete("/packages/{id}", packageHandler.Delete)

			r.Get("/auth/validate", authHandler.ValidateToken)
		})
	})

	// Mount sub-routers on main router
	mainRouter.Mount("/", monitorRouter) // Monitoring endpoints at root level
	mainRouter.Mount("/api", r)          // API endpoints under /api

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, mainRouter); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
