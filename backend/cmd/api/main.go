package main

import (
	"fmt"
	"log"
	"net/http"

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

	// Initialize repositories
	bookingRepo := database.NewBookingRepository(db)

	// Initialize handlers
	bookingHandler := handlers.NewBookingHandler(bookingRepo)

	// Initialize router
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(custommiddleware.CORS(cfg.AllowOrigins))

	// Routes
	r.Route("/api/v1", func(r chi.Router) {
		// Bookings
		r.Route("/bookings", func(r chi.Router) {
			r.Get("/", bookingHandler.GetAll)
			r.Post("/", bookingHandler.Create)
			r.Get("/{id}", bookingHandler.GetByID)
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
