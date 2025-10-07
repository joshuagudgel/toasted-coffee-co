package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
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
	r := chi.NewRouter()

	// Global middleware
	r.Use(custommiddleware.SecureHTTPS)
	r.Use(custommiddleware.SecurityHeaders)
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
			r.Get("/packages", packageHandler.GetAll)
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

			// Package management
			r.Post("/packages", packageHandler.Create)
			r.Get("/packages/{id}", packageHandler.GetByID)
			r.Put("/packages/{id}", packageHandler.Update)
			r.Delete("/packages/{id}", packageHandler.Delete)

			// Auth validation
			r.Get("/auth/validate", authHandler.ValidateToken)
		})
	})

	// Health check with minimal rate limiting
	r.With(httprate.LimitByIP(10, 1*time.Minute)).Get("/health", func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		healthCheckCount++

		// Log detailed request information
		log.Printf("HEALTH CHECK #%d: Started at %v", healthCheckCount, startTime)
		log.Printf("HEALTH CHECK: Method=%s, URL=%s, RemoteAddr=%s, UserAgent=%s",
			r.Method, r.URL.String(), r.RemoteAddr, r.Header.Get("User-Agent"))
		log.Printf("HEALTH CHECK: Headers: %+v", r.Header)

		// Detect if this might be a wake-up after sleep
		if !lastHealthCheck.IsZero() {
			timeSinceLastCheck := startTime.Sub(lastHealthCheck)
			if timeSinceLastCheck > 2*time.Minute {
				log.Printf("HEALTH CHECK: POTENTIAL WAKE-UP DETECTED - Last check was %v ago", timeSinceLastCheck)
			}
		}
		lastHealthCheck = startTime

		// Prepare response data
		response := map[string]interface{}{
			"status":         "healthy",
			"timestamp":      startTime.Format(time.RFC3339),
			"uptime":         startTime.Sub(serviceStartTime).String(),
			"check_count":    healthCheckCount,
			"environment":    os.Getenv("ENVIRONMENT"),
			"go_version":     runtime.Version(),
			"num_goroutines": runtime.NumGoroutine(),
		}

		// Memory stats
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		response["memory"] = map[string]interface{}{
			"alloc_mb":       memStats.Alloc / 1024 / 1024,
			"total_alloc_mb": memStats.TotalAlloc / 1024 / 1024,
			"sys_mb":         memStats.Sys / 1024 / 1024,
			"num_gc":         memStats.NumGC,
		}

		// Database health check with timeout
		dbCtx, dbCancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer dbCancel()

		dbCheckStart := time.Now()
		if err := db.Pool.Ping(dbCtx); err != nil {
			log.Printf("HEALTH CHECK: Database ping FAILED: %v", err)
			response["database"] = map[string]interface{}{
				"status":        "failed",
				"error":         err.Error(),
				"response_time": time.Since(dbCheckStart).String(),
			}
			response["status"] = "unhealthy"

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)

			if jsonErr := json.NewEncoder(w).Encode(response); jsonErr != nil {
				log.Printf("HEALTH CHECK: Failed to encode JSON response: %v", jsonErr)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			log.Printf("HEALTH CHECK: FAILED - Total time: %v", time.Since(startTime))
			return
		}

		// Test database query
		var dbResult int
		queryStart := time.Now()
		if err := db.Pool.QueryRow(dbCtx, "SELECT 1").Scan(&dbResult); err != nil {
			log.Printf("HEALTH CHECK: Database query FAILED: %v", err)
			response["database"] = map[string]interface{}{
				"status":     "query_failed",
				"error":      err.Error(),
				"ping_time":  time.Since(dbCheckStart).String(),
				"query_time": time.Since(queryStart).String(),
			}
			response["status"] = "unhealthy"

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(response)

			log.Printf("HEALTH CHECK: FAILED - Total time: %v", time.Since(startTime))
			return
		}

		// Database is healthy
		response["database"] = map[string]interface{}{
			"status":       "healthy",
			"ping_time":    time.Since(dbCheckStart).String(),
			"query_time":   time.Since(queryStart).String(),
			"query_result": dbResult,
		}

		// Add database connection pool stats if available
		if stats := db.Pool.Stat(); stats != nil {
			response["database_pool"] = map[string]interface{}{
				"total_conns":    stats.TotalConns(),
				"acquired_conns": stats.AcquiredConns(),
				"idle_conns":     stats.IdleConns(),
				"max_conns":      stats.MaxConns(),
			}
		}

		// Add request processing time
		response["response_time"] = time.Since(startTime).String()

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Health-Check-Count", fmt.Sprintf("%d", healthCheckCount))
		w.Header().Set("X-Uptime", startTime.Sub(serviceStartTime).String())
		w.WriteHeader(http.StatusOK)

		// Encode and send response
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("HEALTH CHECK: Failed to encode successful response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}

		log.Printf("HEALTH CHECK: SUCCESS - Total time: %v", time.Since(startTime))
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
