package database_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/database"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/models"
)

// TestDB wraps pgxpool for testing
type TestDB struct {
	Pool *pgxpool.Pool
}

func setupTestDB(t *testing.T) *TestDB {
	// Get test database URL from environment or use default
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/toasted_coffee_test?sslmode=disable"
	}

	// Connect to database
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Create tables
	_, err = pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS bookings (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            email VARCHAR(255),
            phone VARCHAR(20),
            date DATE NOT NULL,
            time VARCHAR(10) NOT NULL,
            people INTEGER NOT NULL,
            location VARCHAR(255) NOT NULL,
            notes TEXT,
            coffee_flavors VARCHAR[] NOT NULL,
            milk_options VARCHAR[] NOT NULL,
            package VARCHAR(100),
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	return &TestDB{Pool: pool}
}

func cleanupTestDB(t *testing.T, db *TestDB) {
	// Clean up test data
	_, err := db.Pool.Exec(context.Background(), "DELETE FROM bookings")
	if err != nil {
		t.Fatalf("Failed to clean up test database: %v", err)
	}
	db.Pool.Close()
}

func TestCreateBooking(t *testing.T) {
	// Skip test if no database is available
	if os.Getenv("SKIP_DB_TESTS") == "true" {
		t.Skip("Skipping database tests")
	}

	testDB := setupTestDB(t)
	defer cleanupTestDB(t, testDB)

	// Create wrapped DB object
	db := &database.DB{Pool: testDB.Pool}
	repo := database.NewBookingRepository(db)

	tests := []struct {
		name        string
		booking     *models.Booking
		expectError bool
	}{
		{
			name: "Valid booking with email",
			booking: &models.Booking{
				Name:          "Test User",
				Email:         "test@example.com",
				Date:          "2025-06-01",
				Time:          "14:00",
				People:        5,
				Location:      "Test Location",
				CoffeeFlavors: []string{"french_toast"},
				MilkOptions:   []string{"whole"},
				Package:       "Group",
			},
			expectError: false,
		},
		{
			name: "Valid booking with phone",
			booking: &models.Booking{
				Name:          "Test User 2",
				Phone:         "555-1234",
				Date:          "2025-06-01",
				Time:          "15:00",
				People:        10,
				Location:      "Another Location",
				CoffeeFlavors: []string{"mexican_mocha"},
				MilkOptions:   []string{"oat"},
			},
			expectError: false,
		},
		{
			name: "Invalid date format",
			booking: &models.Booking{
				Name:          "Test User 3",
				Email:         "test3@example.com",
				Date:          "invalid-date", // Invalid date
				Time:          "16:00",
				People:        15,
				Location:      "Test Location 3",
				CoffeeFlavors: []string{"cinnamon_brown_sugar"},
				MilkOptions:   []string{"almond"},
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create booking
			id, err := repo.Create(context.Background(), tc.booking)

			// Check for expected errors
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			// If not expecting error, but got one
			if err != nil {
				t.Fatalf("Failed to create booking: %v", err)
			}

			// Check that ID was returned
			if id <= 0 {
				t.Errorf("Expected positive ID, got %d", id)
			}

			// Verify booking was saved correctly
			savedBooking, err := repo.GetByID(context.Background(), id)
			if err != nil {
				t.Fatalf("Failed to retrieve booking: %v", err)
			}

			// Verify fields were saved correctly
			if savedBooking.Name != tc.booking.Name {
				t.Errorf("Expected name %s, got %s", tc.booking.Name, savedBooking.Name)
			}

			if savedBooking.Email != tc.booking.Email {
				t.Errorf("Expected email %s, got %s", tc.booking.Email, savedBooking.Email)
			}

			if savedBooking.Phone != tc.booking.Phone {
				t.Errorf("Expected phone %s, got %s", tc.booking.Phone, savedBooking.Phone)
			}

			if savedBooking.Package != tc.booking.Package {
				t.Errorf("Expected package %s, got %s", tc.booking.Package, savedBooking.Package)
			}
		})
	}
}

func TestGetAllBookings_EdgeCases(t *testing.T) {
	// Skip test if no database is available
	if os.Getenv("SKIP_DB_TESTS") == "true" {
		t.Skip("Skipping database tests")
	}

	testDB := setupTestDB(t)
	defer cleanupTestDB(t, testDB)

	// Create wrapped DB object
	db := &database.DB{Pool: testDB.Pool}
	repo := database.NewBookingRepository(db)

	// Test 1: Empty database
	t.Run("Empty database", func(t *testing.T) {
		bookings, err := repo.GetAll(context.Background())
		if err != nil {
			t.Fatalf("Failed to retrieve bookings: %v", err)
		}

		if len(bookings) != 0 {
			t.Errorf("Expected empty result, got %d bookings", len(bookings))
		}
	})

	// Test 2: Insert multiple bookings
	t.Run("Multiple bookings", func(t *testing.T) {
		// Insert test bookings
		for i := 0; i < 5; i++ {
			booking := &models.Booking{
				Name:          fmt.Sprintf("User %d", i),
				Email:         fmt.Sprintf("user%d@example.com", i),
				Date:          "2025-06-01",
				Time:          "14:00",
				People:        5,
				Location:      "Test Location",
				CoffeeFlavors: []string{"french_toast"},
				MilkOptions:   []string{"whole"},
			}
			_, err := repo.Create(context.Background(), booking)
			if err != nil {
				t.Fatalf("Failed to create test booking: %v", err)
			}
		}

		// Retrieve all bookings
		bookings, err := repo.GetAll(context.Background())
		if err != nil {
			t.Fatalf("Failed to retrieve bookings: %v", err)
		}

		if len(bookings) != 5 {
			t.Errorf("Expected 5 bookings, got %d", len(bookings))
		}
	})
}
