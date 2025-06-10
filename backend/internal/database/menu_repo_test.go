package database_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/database"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/models"
)

// setupMenuTestDB initializes a test database for menu items tests
func setupMenuTestDB(t *testing.T) *TestDB {
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
    CREATE TABLE IF NOT EXISTS menu_items (
        id SERIAL PRIMARY KEY,
        value VARCHAR(100) NOT NULL,
        label VARCHAR(100) NOT NULL,
        type VARCHAR(20) NOT NULL,
        active BOOLEAN DEFAULT true,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )
    `)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	return &TestDB{Pool: pool}
}

// cleanupMenuTestDB cleans up the test database
func cleanupMenuTestDB(t *testing.T, db *TestDB) {
	// Clean up test data
	_, err := db.Pool.Exec(context.Background(), "DELETE FROM menu_items")
	if err != nil {
		t.Fatalf("Failed to clean up test database: %v", err)
	}
	db.Pool.Close()
}

func TestCreateMenuItem(t *testing.T) {
	log.Println("Running TestCreateMenuItem...")
	// Skip test if no database is available
	if os.Getenv("SKIP_DB_TESTS") == "true" {
		t.Skip("Skipping database tests")
	}

	testDB := setupMenuTestDB(t)
	defer cleanupMenuTestDB(t, testDB)

	// Create wrapped DB object
	db := &database.DB{Pool: testDB.Pool}
	repo := database.NewMenuRepository(db)

	tests := []struct {
		name        string
		menuItem    *models.MenuItem
		expectError bool
	}{
		{
			name: "Valid coffee flavor",
			menuItem: &models.MenuItem{
				Value:  "french_toast",
				Label:  "French Toast",
				Type:   models.CoffeeFlavor,
				Active: true,
			},
			expectError: false,
		},
		{
			name: "Valid milk option",
			menuItem: &models.MenuItem{
				Value:  "whole",
				Label:  "Whole Milk",
				Type:   models.MilkOption,
				Active: true,
			},
			expectError: false,
		},
		{
			name: "Inactive coffee flavor",
			menuItem: &models.MenuItem{
				Value:  "seasonal_pumpkin",
				Label:  "Seasonal Pumpkin",
				Type:   models.CoffeeFlavor,
				Active: false,
			},
			expectError: false,
		},
		{
			name: "Duplicate value", // This may cause an error if database has a unique constraint
			menuItem: &models.MenuItem{
				Value:  "french_toast", // Same as the first test case
				Label:  "French Toast Duplicate",
				Type:   models.CoffeeFlavor,
				Active: true,
			},
			expectError: false, // Change to true if your database has a unique constraint
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create menu item
			id, err := repo.Create(context.Background(), tc.menuItem)

			// Check for expected errors
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			// If not expecting error, but got one
			if err != nil {
				t.Fatalf("Failed to create menu item: %v", err)
			}

			// Check that ID was returned
			if id <= 0 {
				t.Errorf("Expected positive ID, got %d", id)
			}

			// Retrieve all items to verify creation
			items, err := repo.GetAll(context.Background())
			if err != nil {
				t.Fatalf("Failed to retrieve menu items: %v", err)
			}

			// Find our item in the list
			var found bool
			for _, item := range items {
				if item.ID == id {
					found = true

					// Verify fields were saved correctly
					if item.Value != tc.menuItem.Value {
						t.Errorf("Expected value %s, got %s", tc.menuItem.Value, item.Value)
					}

					if item.Label != tc.menuItem.Label {
						t.Errorf("Expected label %s, got %s", tc.menuItem.Label, item.Label)
					}

					if item.Type != tc.menuItem.Type {
						t.Errorf("Expected type %s, got %s", tc.menuItem.Type, item.Type)
					}

					if item.Active != tc.menuItem.Active {
						t.Errorf("Expected active %v, got %v", tc.menuItem.Active, item.Active)
					}

					break
				}
			}

			if !found {
				t.Errorf("Newly created menu item with ID %d was not found in GetAll results", id)
			}
		})
	}
}

func TestGetByType(t *testing.T) {
	log.Println("Running TestGetByType...")
	// Skip test if no database is available
	if os.Getenv("SKIP_DB_TESTS") == "true" {
		t.Skip("Skipping database tests")
	}

	testDB := setupMenuTestDB(t)
	defer cleanupMenuTestDB(t, testDB)

	// Create wrapped DB object
	db := &database.DB{Pool: testDB.Pool}
	repo := database.NewMenuRepository(db)

	// Create test data
	testItems := []*models.MenuItem{
		{Value: "french_toast", Label: "French Toast", Type: models.CoffeeFlavor, Active: true},
		{Value: "mexican_mocha", Label: "Mexican Mocha", Type: models.CoffeeFlavor, Active: true},
		{Value: "seasonal_pumpkin", Label: "Seasonal Pumpkin", Type: models.CoffeeFlavor, Active: false},
		{Value: "whole", Label: "Whole Milk", Type: models.MilkOption, Active: true},
		{Value: "oat", Label: "Oat Milk", Type: models.MilkOption, Active: true},
		{Value: "almond", Label: "Almond Milk", Type: models.MilkOption, Active: false},
	}

	for _, item := range testItems {
		_, err := repo.Create(context.Background(), item)
		if err != nil {
			t.Fatalf("Failed to create test data: %v", err)
		}
	}

	// Test getting coffee flavors
	t.Run("Get coffee flavors", func(t *testing.T) {
		items, err := repo.GetByType(context.Background(), models.CoffeeFlavor)
		if err != nil {
			t.Fatalf("Failed to get coffee flavors: %v", err)
		}

		// Check count (should be 3 coffee flavors)
		if len(items) != 3 {
			t.Errorf("Expected 3 coffee flavors, got %d", len(items))
		}

		// Verify all items are coffee flavors
		for _, item := range items {
			if item.Type != models.CoffeeFlavor {
				t.Errorf("Expected type %s, got %s", models.CoffeeFlavor, item.Type)
			}
		}
	})

	// Test getting milk options
	t.Run("Get milk options", func(t *testing.T) {
		items, err := repo.GetByType(context.Background(), models.MilkOption)
		if err != nil {
			t.Fatalf("Failed to get milk options: %v", err)
		}

		// Check count (should be 3 milk options)
		if len(items) != 3 {
			t.Errorf("Expected 3 milk options, got %d", len(items))
		}

		// Verify all items are milk options
		for _, item := range items {
			if item.Type != models.MilkOption {
				t.Errorf("Expected type %s, got %s", models.MilkOption, item.Type)
			}
		}
	})
}

func TestUpdateMenuItem(t *testing.T) {
	log.Println("Running TestUpdateMenuItem...")
	// Skip test if no database is available
	if os.Getenv("SKIP_DB_TESTS") == "true" {
		t.Skip("Skipping database tests")
	}

	testDB := setupMenuTestDB(t)
	defer cleanupMenuTestDB(t, testDB)

	// Create wrapped DB object
	db := &database.DB{Pool: testDB.Pool}
	repo := database.NewMenuRepository(db)

	// Create a test item
	item := &models.MenuItem{
		Value:  "vanilla_chai",
		Label:  "Vanilla Chai",
		Type:   models.CoffeeFlavor,
		Active: true,
	}

	id, err := repo.Create(context.Background(), item)
	if err != nil {
		t.Fatalf("Failed to create test item: %v", err)
	}

	// Update the item
	updatedItem := &models.MenuItem{
		Value:  "dirty_vanilla_chai",
		Label:  "Dirty Vanilla Chai",
		Type:   models.CoffeeFlavor,
		Active: false,
	}

	err = repo.Update(context.Background(), id, updatedItem)
	if err != nil {
		t.Fatalf("Failed to update menu item: %v", err)
	}

	// Retrieve all items to verify update
	items, err := repo.GetAll(context.Background())
	if err != nil {
		t.Fatalf("Failed to retrieve menu items: %v", err)
	}

	// Find our updated item
	var found bool
	for _, item := range items {
		if item.ID == id {
			found = true

			// Verify fields were updated correctly
			if item.Value != updatedItem.Value {
				t.Errorf("Expected value %s, got %s", updatedItem.Value, item.Value)
			}

			if item.Label != updatedItem.Label {
				t.Errorf("Expected label %s, got %s", updatedItem.Label, item.Label)
			}

			if item.Type != updatedItem.Type {
				t.Errorf("Expected type %s, got %s", updatedItem.Type, item.Type)
			}

			if item.Active != updatedItem.Active {
				t.Errorf("Expected active %v, got %v", updatedItem.Active, item.Active)
			}

			break
		}
	}

	if !found {
		t.Error("Updated menu item was not found in GetAll results")
	}

	// Test updating non-existent item
	err = repo.Update(context.Background(), 9999, updatedItem)
	if err == nil {
		t.Error("Expected error when updating non-existent item, but got nil")
	}
}

func TestDeleteMenuItem(t *testing.T) {
	log.Println("Running TestDeleteMenuItem...")
	// Skip test if no database is available
	if os.Getenv("SKIP_DB_TESTS") == "true" {
		t.Skip("Skipping database tests")
	}

	testDB := setupMenuTestDB(t)
	defer cleanupMenuTestDB(t, testDB)

	// Create wrapped DB object
	db := &database.DB{Pool: testDB.Pool}
	repo := database.NewMenuRepository(db)

	// Create a test item
	item := &models.MenuItem{
		Value:  "to_be_deleted",
		Label:  "To Be Deleted",
		Type:   models.CoffeeFlavor,
		Active: true,
	}

	id, err := repo.Create(context.Background(), item)
	if err != nil {
		t.Fatalf("Failed to create test item: %v", err)
	}

	// Verify item exists
	items, err := repo.GetAll(context.Background())
	if err != nil {
		t.Fatalf("Failed to retrieve menu items: %v", err)
	}

	var found bool
	for _, item := range items {
		if item.ID == id {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("Test item not found before deletion")
	}

	// Delete the item
	err = repo.Delete(context.Background(), id)
	if err != nil {
		t.Fatalf("Failed to delete menu item: %v", err)
	}

	// Verify item is gone
	items, err = repo.GetAll(context.Background())
	if err != nil {
		t.Fatalf("Failed to retrieve menu items after deletion: %v", err)
	}

	for _, item := range items {
		if item.ID == id {
			t.Error("Deleted item still exists in database")
		}
	}

	// Test deleting non-existent item
	err = repo.Delete(context.Background(), 9999)
	if err == nil {
		t.Error("Expected error when deleting non-existent item, but got nil")
	}
}

func TestGetAllMenuItems(t *testing.T) {
	log.Println("Running TestGetAllMenuItems...")
	// Skip test if no database is available
	if os.Getenv("SKIP_DB_TESTS") == "true" {
		t.Skip("Skipping database tests")
	}

	testDB := setupMenuTestDB(t)
	defer cleanupMenuTestDB(t, testDB)

	// Create wrapped DB object
	db := &database.DB{Pool: testDB.Pool}
	repo := database.NewMenuRepository(db)

	// Test with no items
	t.Run("Empty database", func(t *testing.T) {
		items, err := repo.GetAll(context.Background())
		if err != nil {
			t.Fatalf("Failed to get menu items: %v", err)
		}

		if len(items) != 0 {
			t.Errorf("Expected 0 items, got %d", len(items))
		}
	})

	// Create test items
	testItems := []*models.MenuItem{
		{Value: "item1", Label: "Item 1", Type: models.CoffeeFlavor, Active: true},
		{Value: "item2", Label: "Item 2", Type: models.CoffeeFlavor, Active: false},
		{Value: "item3", Label: "Item 3", Type: models.MilkOption, Active: true},
		{Value: "item4", Label: "Item 4", Type: models.MilkOption, Active: false},
	}

	// Insert test items
	for _, item := range testItems {
		_, err := repo.Create(context.Background(), item)
		if err != nil {
			t.Fatalf("Failed to create test data: %v", err)
		}

		// Small sleep to ensure timestamps are different
		time.Sleep(10 * time.Millisecond)
	}

	// Test with items
	t.Run("Multiple items", func(t *testing.T) {
		items, err := repo.GetAll(context.Background())
		if err != nil {
			t.Fatalf("Failed to get menu items: %v", err)
		}

		// Should have all test items
		if len(items) != len(testItems) {
			t.Errorf("Expected %d items, got %d", len(testItems), len(items))
		}

		// Verify we have correct counts for each type
		coffeeCount := 0
		milkCount := 0
		activeCount := 0
		inactiveCount := 0

		for _, item := range items {
			if item.Type == models.CoffeeFlavor {
				coffeeCount++
			} else if item.Type == models.MilkOption {
				milkCount++
			}

			if item.Active {
				activeCount++
			} else {
				inactiveCount++
			}
		}

		if coffeeCount != 2 {
			t.Errorf("Expected 2 coffee flavors, got %d", coffeeCount)
		}

		if milkCount != 2 {
			t.Errorf("Expected 2 milk options, got %d", milkCount)
		}

		if activeCount != 2 {
			t.Errorf("Expected 2 active items, got %d", activeCount)
		}

		if inactiveCount != 2 {
			t.Errorf("Expected 2 inactive items, got %d", inactiveCount)
		}
	})
}

func TestActiveItems(t *testing.T) {
	log.Println("Running TestActiveItems...")
	// Skip test if no database is available
	if os.Getenv("SKIP_DB_TESTS") == "true" {
		t.Skip("Skipping database tests")
	}

	testDB := setupMenuTestDB(t)
	defer cleanupMenuTestDB(t, testDB)

	// Create wrapped DB object
	db := &database.DB{Pool: testDB.Pool}
	repo := database.NewMenuRepository(db)

	// Create test items with different active states
	testItems := []*models.MenuItem{
		{Value: "active_coffee", Label: "Active Coffee", Type: models.CoffeeFlavor, Active: true},
		{Value: "inactive_coffee", Label: "Inactive Coffee", Type: models.CoffeeFlavor, Active: false},
		{Value: "active_milk", Label: "Active Milk", Type: models.MilkOption, Active: true},
		{Value: "inactive_milk", Label: "Inactive Milk", Type: models.MilkOption, Active: false},
	}

	for _, item := range testItems {
		_, err := repo.Create(context.Background(), item)
		if err != nil {
			t.Fatalf("Failed to create test data: %v", err)
		}
	}

	// Test toggling active state
	t.Run("Toggle active state", func(t *testing.T) {
		// Get all items to find IDs
		items, err := repo.GetAll(context.Background())
		if err != nil {
			t.Fatalf("Failed to get menu items: %v", err)
		}

		var inactiveItemID int
		var activeItemID int

		// Find an active and inactive item
		for _, item := range items {
			if item.Active && activeItemID == 0 {
				activeItemID = item.ID
			} else if !item.Active && inactiveItemID == 0 {
				inactiveItemID = item.ID
			}

			if activeItemID != 0 && inactiveItemID != 0 {
				break
			}
		}

		// Make active item inactive
		if activeItemID != 0 {
			updateItem := &models.MenuItem{
				Value:  "updated_name",
				Label:  "Updated Name",
				Type:   models.CoffeeFlavor,
				Active: false, // Change to inactive
			}
			err = repo.Update(context.Background(), activeItemID, updateItem)
			if err != nil {
				t.Fatalf("Failed to update active item: %v", err)
			}
		}

		// Make inactive item active
		if inactiveItemID != 0 {
			updateItem := &models.MenuItem{
				Value:  "updated_inactive",
				Label:  "Updated Inactive",
				Type:   models.MilkOption,
				Active: true, // Change to active
			}
			err = repo.Update(context.Background(), inactiveItemID, updateItem)
			if err != nil {
				t.Fatalf("Failed to update inactive item: %v", err)
			}
		}

		// Get all items again to verify changes
		updatedItems, err := repo.GetAll(context.Background())
		if err != nil {
			t.Fatalf("Failed to get menu items after update: %v", err)
		}

		// Check that active states were toggled
		for _, item := range updatedItems {
			if item.ID == activeItemID {
				if item.Active {
					t.Error("Item should be inactive after update")
				}
				if item.Value != "updated_name" {
					t.Errorf("Expected value 'updated_name', got '%s'", item.Value)
				}
			}

			if item.ID == inactiveItemID {
				if !item.Active {
					t.Error("Item should be active after update")
				}
				if item.Value != "updated_inactive" {
					t.Errorf("Expected value 'updated_inactive', got '%s'", item.Value)
				}
			}
		}
	})
}
