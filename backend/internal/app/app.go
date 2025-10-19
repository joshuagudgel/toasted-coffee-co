package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joshuagudgel/toasted-coffee/backend/internal/config"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/database"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/handlers"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/server"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/services"
)

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

	// Run migrations - Admin seeder and other migrations
	if err := runDatabaseSetup(db); err != nil {
		db.Close()
		return nil, err
	}

	// Initialize services
	emailService := services.NewEmailService()

	// Initialize repositories
	repos := database.NewRepositories(db)

	// Initialize handlers
	handlers := handlers.NewHandlers(repos, emailService)

	// Setup router
	router := server.NewRouter(handlers, cfg)

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

func runDatabaseSetup(db *database.DB) error {
	migrator := database.NewMigrator(db)
	if err := migrator.RunMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	seeder := database.NewSeeder(db)
	if err := seeder.SeedAdminUser(); err != nil {
		return fmt.Errorf("failed to seed admin user: %w", err)
	}

	return nil
}

func (a *App) Close() error {
	if a.db != nil {
		a.db.Close()
	}
	return nil
}
