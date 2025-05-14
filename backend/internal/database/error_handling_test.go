package database_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/database"
)

func TestDatabaseConnectionHandling(t *testing.T) {
	// Test with invalid connection string
	t.Run("Invalid connection string", func(t *testing.T) {
		_, err := database.New("postgres://invalid:invalid@nonexistent:5432/db")
		if err == nil {
			t.Error("Expected error with invalid connection string, got nil")
		}
	})

	// Test connection timeout handling
	t.Run("Connection timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		// This should time out quickly
		config, _ := pgxpool.ParseConfig("postgres://postgres:postgres@localhost:5432/toasted_coffee_test")
		_, err := pgxpool.NewWithConfig(ctx, config)
		if err == nil {
			t.Error("Expected timeout error, got nil")
		}
	})
}

func TestTransactionHandling(t *testing.T) {
	// Skip test if no database is available
	if os.Getenv("SKIP_DB_TESTS") == "true" {
		t.Skip("Skipping database transaction tests")
	}

	testDB := setupTestDB(t)
	defer cleanupTestDB(t, testDB)

	t.Run("Transaction rollback on error", func(t *testing.T) {
		// Start a transaction
		tx, err := testDB.Pool.Begin(context.Background())
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}

		// Insert a test booking
		_, err = tx.Exec(context.Background(), `
            INSERT INTO bookings (name, date, time, people, location, coffee_flavors, milk_options)
            VALUES ($1, $2, $3, $4, $5, $6, $7)
        `, "Transaction Test", "2025-06-01", "14:00", 5, "Test Location", []string{"french_toast"}, []string{"whole"})
		if err != nil {
			t.Fatalf("Failed to insert test booking: %v", err)
		}

		// Rollback the transaction
		err = tx.Rollback(context.Background())
		if err != nil {
			t.Fatalf("Failed to rollback transaction: %v", err)
		}

		// Verify booking was not inserted
		var count int
		err = testDB.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM bookings WHERE name = $1", "Transaction Test").Scan(&count)
		if err != nil {
			t.Fatalf("Failed to query bookings: %v", err)
		}

		if count != 0 {
			t.Errorf("Expected 0 bookings after rollback, got %d", count)
		}
	})
}
