package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
)

const (
	PublicReadLimit  = 100
	PublicWriteLimit = 10
	ContactLimit     = 5
	AuthLimit        = 20
	AdminLimit       = 200
	HealthCheckLimit = 60
)

var serviceStartTime = time.Now()

type App struct {
	cfg    *config.Config
	db     *database.DB
	server *http.Server
}

func New() (*App, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Connect to database
	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations
	migrator := database.NewMigrator(db)
	if err := migrator.RunMigrations(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Seed admin user
	seeder := database.NewSeeder(db)
	if err := seeder.SeedAdminUser(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to seed admin user: %w", err)
	}

	// Initialize services
	emailService := services.NewEmailService()

	// Initialize repositories
	repos := database.NewRepositories(db)

	// Initialize handlers
	handlers := handlers.NewHandlers(repos, emailService)

	// Setup router
	router := setupRouter(handlers, cfg)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	return &App{
		cfg:    cfg,
		db:     db,
		server: httpServer,
	}, nil
}

func (a *App) Run() error {
	log.Printf("Server starting on %s", a.server.Addr)
	return a.server.ListenAndServe()
}

func (a *App) Close() error {
	if a.db != nil {
		a.db.Close()
	}
	return nil
}

func setupRouter(h *handlers.Handlers, cfg *config.Config) *chi.Mux {
	mainRouter := chi.NewRouter()

	// Monitoring sub-router
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

	monitorRouter.Get("/ping-simple", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	monitorRouter.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		requestTime := time.Now()
		log.Printf("PING REQUEST: time=%v, ip=%s, user_agent=%s",
			requestTime.Format(time.RFC3339),
			r.RemoteAddr,
			r.Header.Get("User-Agent"))

		userAgent := r.Header.Get("User-Agent")
		if strings.Contains(strings.ToLower(userAgent), "cron") {
			log.Printf("PING: CRON JOB DETECTED - UserAgent: %s", userAgent)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":    "ok",
			"timestamp": requestTime.Format(time.RFC3339),
		})

		log.Printf("PING SUCCESS: Response sent in %v", time.Since(requestTime))
	})

	monitorRouter.Get("/test-render", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("RENDER TEST: Request received at %v", time.Now())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		log.Printf("RENDER TEST: Response sent")
	})

	// API sub-router
	apiRouter := chi.NewRouter()
	apiRouter.Use(custommiddleware.SecureHTTPS)
	apiRouter.Use(custommiddleware.SecurityHeaders)
	apiRouter.Use(middleware.Logger)
	apiRouter.Use(middleware.Recoverer)
	apiRouter.Use(custommiddleware.CORS(cfg.AllowOrigins))

	apiRouter.Route("/v1", func(r chi.Router) {
		// Public read-only endpoints
		r.Group(func(r chi.Router) {
			r.Use(httprate.LimitByIP(PublicReadLimit, 1*time.Minute))
			r.Get("/menu", h.Menu.GetAll)
			r.Get("/menu/{type}", h.Menu.GetByType)
			r.Get("/packages", h.Package.GetAll)
		})

		// Public write endpoints
		r.Group(func(r chi.Router) {
			r.Use(httprate.LimitByIP(PublicWriteLimit, 1*time.Minute))
			r.Post("/bookings", h.Booking.Create)
		})

		// Contact endpoint
		r.With(httprate.LimitByIP(ContactLimit, 1*time.Minute)).
			Post("/contact", h.Contact.HandleInquiry)

		// Authentication endpoints
		r.Group(func(r chi.Router) {
			r.Use(httprate.LimitByIP(AuthLimit, 1*time.Minute))
			r.Post("/auth/login", h.Auth.Login)
			r.Post("/auth/refresh", h.Auth.RefreshToken)
			r.Post("/auth/logout", h.Auth.Logout)
		})

		// Protected admin endpoints
		r.Group(func(r chi.Router) {
			r.Use(custommiddleware.JWTAuth)
			r.Use(httprate.LimitByIP(AdminLimit, 1*time.Minute))

			r.Get("/bookings", h.Booking.GetAll)
			r.Get("/bookings/{id}", h.Booking.GetByID)
			r.Put("/bookings/{id}", h.Booking.Update)
			r.Delete("/bookings/{id}", h.Booking.Delete)
			r.Post("/bookings/{id}/archive", h.Booking.Archive)
			r.Post("/bookings/{id}/unarchive", h.Booking.Unarchive)

			r.Post("/menu", h.Menu.Create)
			r.Put("/menu/{id}", h.Menu.Update)
			r.Delete("/menu/{id}", h.Menu.Delete)

			r.Post("/packages", h.Package.Create)
			r.Get("/packages/{id}", h.Package.GetByID)
			r.Put("/packages/{id}", h.Package.Update)
			r.Delete("/packages/{id}", h.Package.Delete)

			r.Get("/auth/validate", h.Auth.ValidateToken)
		})
	})

	mainRouter.Mount("/", monitorRouter)
	mainRouter.Mount("/api", apiRouter)
	return mainRouter
}
