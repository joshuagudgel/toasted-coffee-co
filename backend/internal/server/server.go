package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/config"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/handlers"
	custommiddleware "github.com/joshuagudgel/toasted-coffee/backend/internal/middleware"
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

func NewRouter(h *handlers.Handlers, cfg *config.Config) *chi.Mux {
	mainRouter := chi.NewRouter()

	// Mount sub-routers for better organization
	mainRouter.Mount("/", newMonitorRouter())
	mainRouter.Mount("/api", newAPIRouter(h, cfg))

	return mainRouter
}

func newMonitorRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/health", healthHandler)
	router.Get("/ping-simple", pingSimpleHandler)
	router.Get("/ping", pingHandler)
	router.Get("/test-render", testRenderHandler)

	return router
}

func newAPIRouter(h *handlers.Handlers, cfg *config.Config) *chi.Mux {
	router := chi.NewRouter()

	// Common middleware
	router.Use(custommiddleware.SecureHTTPS)
	router.Use(custommiddleware.SecurityHeaders)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(custommiddleware.CORS(cfg.AllowOrigins))

	router.Route("/v1", func(r chi.Router) {
		setupPublicRoutes(r, h)
		setupAuthRoutes(r, h)
		setupAdminRoutes(r, h)
	})

	return router
}

func setupPublicRoutes(r chi.Router, h *handlers.Handlers) {
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
}

func setupAuthRoutes(r chi.Router, h *handlers.Handlers) {
	r.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(AuthLimit, 1*time.Minute))
		r.Post("/auth/login", h.Auth.Login)
		r.Post("/auth/refresh", h.Auth.RefreshToken)
		r.Post("/auth/logout", h.Auth.Logout)
	})
}

func setupAdminRoutes(r chi.Router, h *handlers.Handlers) {
	r.Group(func(r chi.Router) {
		r.Use(custommiddleware.JWTAuth)
		r.Use(httprate.LimitByIP(AdminLimit, 1*time.Minute))

		// Booking routes
		r.Get("/bookings", h.Booking.GetAll)
		r.Get("/bookings/{id}", h.Booking.GetByID)
		r.Put("/bookings/{id}", h.Booking.Update)
		r.Delete("/bookings/{id}", h.Booking.Delete)
		r.Post("/bookings/{id}/archive", h.Booking.Archive)
		r.Post("/bookings/{id}/unarchive", h.Booking.Unarchive)

		// Menu routes
		r.Post("/menu", h.Menu.Create)
		r.Put("/menu/{id}", h.Menu.Update)
		r.Delete("/menu/{id}", h.Menu.Delete)

		// Package routes
		r.Post("/packages", h.Package.Create)
		r.Get("/packages/{id}", h.Package.GetByID)
		r.Put("/packages/{id}", h.Package.Update)
		r.Delete("/packages/{id}", h.Package.Delete)

		// Auth validation
		r.Get("/auth/validate", h.Auth.ValidateToken)
	})
}

// Monitor handler functions
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    time.Since(serviceStartTime).String(),
	})
}

func pingSimpleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
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
}

func testRenderHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("RENDER TEST: Request received at %v", time.Now())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	log.Printf("RENDER TEST: Response sent")
}
