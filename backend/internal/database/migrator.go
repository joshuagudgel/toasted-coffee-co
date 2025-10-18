package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Migrator struct {
	db *DB
}

func NewMigrator(db *DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) RunMigrations() error {
	log.Println("Running database migrations...")

	migrationFiles, err := m.getMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	for _, file := range migrationFiles {
		if err := m.runMigration(file); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", file, err)
		}
	}

	log.Println("All migrations completed successfully")
	return nil
}

func (m *Migrator) getMigrationFiles() ([]string, error) {
	migrationDir := "internal/database/migrations"

	entries, err := os.ReadDir(migrationDir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}

	sort.Strings(files) // Ensure migrations run in order
	return files, nil
}

func (m *Migrator) runMigration(filename string) error {
	migrationPath := filepath.Join("internal/database/migrations", filename)

	migrationSQL, err := os.ReadFile(migrationPath)
	if err != nil {
		log.Printf("Warning: Could not read migration file %s: %v", filename, err)
		return nil // Non-fatal for missing files
	}

	_, err = m.db.Pool.Exec(context.Background(), string(migrationSQL))
	if err != nil {
		if m.isMigrationAlreadyApplied(err) {
			log.Printf("Migration %s already applied, skipping", filename)
			return nil
		}
		return err
	}

	log.Printf("Migration %s executed successfully", filename)
	return nil
}

func (m *Migrator) isMigrationAlreadyApplied(err error) bool {
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "already exists") ||
		strings.Contains(errMsg, "duplicate column")
}
