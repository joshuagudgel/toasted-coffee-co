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
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        archived BOOLEAN DEFAULT FALSE
    )
	`)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	_, err = pool.Exec(context.Background(), `
    DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT FROM information_schema.columns 
            WHERE table_name = 'bookings' AND column_name = 'archived'
        ) THEN
            ALTER TABLE bookings ADD COLUMN archived BOOLEAN DEFAULT FALSE;
        END IF;
    END
    $$;
`)
	if err != nil {
		t.Fatalf("Failed to add archived column: %v", err)
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
				Date:          "invalid-date",
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

// what were you thinking
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
		bookings, err := repo.GetAll(context.Background(), false)
		if err != nil {
			t.Fatalf("Failed to retrieve bookings: %v", err)
		}

		if len(bookings) != 0 {
			t.Errorf("Expected empty result, got %d bookings", len(bookings))
		}

		// Test 2: Mix of active and archived bookings
		t.Run("Mix of active and archived bookings", func(t *testing.T) {
			// Clear the database first to ensure a clean state
			_, err := testDB.Pool.Exec(context.Background(), "DELETE FROM bookings")
			if err != nil {
				t.Fatalf("Failed to clean test database: %v", err)
			}

			// Insert 8 test bookings - 5 active, 3 archived
			for i := 0; i < 8; i++ {
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
				id, err := repo.Create(context.Background(), booking)
				if err != nil {
					t.Fatalf("Failed to create test booking: %v", err)
				}

				// Archive bookings 5, 6, and 7
				if i >= 5 {
					err = repo.Archive(context.Background(), id)
					if err != nil {
						t.Fatalf("Failed to archive booking: %v", err)
					}
				}
			}

			// Test 2.1: Retrieve active bookings only
			activeBookings, err := repo.GetAll(context.Background(), false)
			if err != nil {
				t.Fatalf("Failed to retrieve active bookings: %v", err)
			}

			if len(activeBookings) != 5 {
				t.Errorf("Expected 5 active bookings, got %d", len(activeBookings))
			}

			// Verify all returned bookings are active (not archived)
			for _, booking := range activeBookings {
				if booking.Archived {
					t.Errorf("GetAll with includeArchived=false returned an archived booking (ID: %d)", booking.ID)
				}
			}

			// Test 2.2: Retrieve all bookings including archived
			allBookings, err := repo.GetAll(context.Background(), true)
			if err != nil {
				t.Fatalf("Failed to retrieve all bookings: %v", err)
			}

			if len(allBookings) != 8 {
				t.Errorf("Expected 8 total bookings, got %d", len(allBookings))
			}

			// Count archived bookings to verify we got the expected number
			archivedCount := 0
			for _, booking := range allBookings {
				if booking.Archived {
					archivedCount++
				}
			}

			if archivedCount != 3 {
				t.Errorf("Expected 3 archived bookings, got %d", archivedCount)
			}
		})
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
		bookings, err := repo.GetAll(context.Background(), false)
		if err != nil {
			t.Fatalf("Failed to retrieve bookings: %v", err)
		}

		if len(bookings) != 5 {
			t.Errorf("Expected 5 bookings, got %d", len(bookings))
		}
	})
}

func TestArchiveAndUnarchiveBooking(t *testing.T) {
	// Skip test if no database is available
	if os.Getenv("SKIP_DB_TESTS") == "true" {
		t.Skip("Skipping database tests")
	}

	testDB := setupTestDB(t)
	defer cleanupTestDB(t, testDB)

	// Create wrapped DB object
	db := &database.DB{Pool: testDB.Pool}
	repo := database.NewBookingRepository(db)

	// Create a test booking
	booking := &models.Booking{
		Name:          "Archive Test User",
		Email:         "archive@test.com",
		Date:          "2025-06-01",
		Time:          "14:00",
		People:        5,
		Location:      "Test Location",
		CoffeeFlavors: []string{"french_toast"},
		MilkOptions:   []string{"whole"},
		Package:       "Group",
	}

	// Step 1: Create the booking
	bookingID, err := repo.Create(context.Background(), booking)
	if err != nil {
		t.Fatalf("Failed to create test booking: %v", err)
	}

	// Step 2: Verify it's not archived by default
	createdBooking, err := repo.GetByID(context.Background(), bookingID)
	if err != nil {
		t.Fatalf("Failed to retrieve booking: %v", err)
	}

	if createdBooking.Archived {
		t.Error("New booking should not be archived by default")
	}

	// Step 3: Archive the booking
	err = repo.Archive(context.Background(), bookingID)
	if err != nil {
		t.Fatalf("Failed to archive booking: %v", err)
	}

	// Step 4: Verify it's now archived
	archivedBooking, err := repo.GetByID(context.Background(), bookingID)
	if err != nil {
		t.Fatalf("Failed to retrieve booking after archiving: %v", err)
	}

	if !archivedBooking.Archived {
		t.Error("Booking should be archived after calling Archive")
	}

	// Step 5: Unarchive the booking
	err = repo.Unarchive(context.Background(), bookingID)
	if err != nil {
		t.Fatalf("Failed to unarchive booking: %v", err)
	}

	// Step 6: Verify it's no longer archived
	unarchivedBooking, err := repo.GetByID(context.Background(), bookingID)
	if err != nil {
		t.Fatalf("Failed to retrieve booking after unarchiving: %v", err)
	}

	if unarchivedBooking.Archived {
		t.Error("Booking should not be archived after calling Unarchive")
	}

	// Test edge cases as subtests
	t.Run("Archive non-existent booking", func(t *testing.T) {
		err := repo.Archive(context.Background(), 99999)
		if err == nil {
			t.Error("Expected error when archiving non-existent booking, got nil")
		}
	})

	t.Run("Unarchive non-existent booking", func(t *testing.T) {
		err := repo.Unarchive(context.Background(), 99999)
		if err == nil {
			t.Error("Expected error when unarchiving non-existent booking, got nil")
		}
	})
}

func TestGetAllWithArchiveFiltering(t *testing.T) {
	// Skip test if no database is available
	if os.Getenv("SKIP_DB_TESTS") == "true" {
		t.Skip("Skipping database tests")
	}

	testDB := setupTestDB(t)
	defer cleanupTestDB(t, testDB)

	// Create wrapped DB object
	db := &database.DB{Pool: testDB.Pool}
	repo := database.NewBookingRepository(db)

	// Clear any existing data
	_, err := testDB.Pool.Exec(context.Background(), "DELETE FROM bookings")
	if err != nil {
		t.Fatalf("Failed to clean test database: %v", err)
	}

	// Insert test bookings (3 active, 2 archived)
	for i := 0; i < 5; i++ {
		booking := &models.Booking{
			Name:          fmt.Sprintf("Filter Test %d", i),
			Email:         fmt.Sprintf("filter%d@test.com", i),
			Date:          "2025-06-01",
			Time:          "14:00",
			People:        5,
			Location:      "Test Location",
			CoffeeFlavors: []string{"french_toast"},
			MilkOptions:   []string{"whole"},
		}

		id, err := repo.Create(context.Background(), booking)
		if err != nil {
			t.Fatalf("Failed to create test booking: %v", err)
		}

		// Archive bookings 3 and 4
		if i >= 3 {
			err = repo.Archive(context.Background(), id)
			if err != nil {
				t.Fatalf("Failed to archive booking: %v", err)
			}
		}
	}

	// Test 1: Get active bookings only
	t.Run("Get active bookings only", func(t *testing.T) {
		activeBookings, err := repo.GetAll(context.Background(), false)
		if err != nil {
			t.Fatalf("Failed to retrieve active bookings: %v", err)
		}

		if len(activeBookings) != 3 {
			t.Errorf("Expected 3 active bookings, got %d", len(activeBookings))
		}

		// Verify none of the returned bookings are archived
		for _, booking := range activeBookings {
			if booking.Archived {
				t.Errorf("GetAll with includeArchived=false returned an archived booking (ID: %d)", booking.ID)
			}
		}
	})

	// Test 2: Get all bookings including archived
	t.Run("Get all bookings including archived", func(t *testing.T) {
		allBookings, err := repo.GetAll(context.Background(), true)
		if err != nil {
			t.Fatalf("Failed to retrieve all bookings: %v", err)
		}

		if len(allBookings) != 5 {
			t.Errorf("Expected 5 total bookings, got %d", len(allBookings))
		}

		// Count archived bookings
		archivedCount := 0
		for _, booking := range allBookings {
			if booking.Archived {
				archivedCount++
			}
		}

		if archivedCount != 2 {
			t.Errorf("Expected 2 archived bookings, got %d", archivedCount)
		}
	})
}
