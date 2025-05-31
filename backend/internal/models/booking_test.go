package models_test

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/models"
)

func TestBookingValidation(t *testing.T) {
	tests := []struct {
		name          string
		booking       models.Booking
		expectedValid bool
		fieldErrors   map[string]bool // Fields expected to have validation errors
	}{
		{
			name: "Fully valid booking",
			booking: models.Booking{
				Name:          "Test User",
				Email:         "test@example.com",
				Date:          "2025-06-01",
				Time:          "14:00",
				People:        5,
				Location:      "Test Location",
				CoffeeFlavors: []string{"french_toast"},
				MilkOptions:   []string{"whole"},
				Package:       "Group",
				Archived:      false,
			},
			expectedValid: true,
		},
		{
			name: "Missing name",
			booking: models.Booking{
				Name:          "",
				Email:         "test@example.com",
				Date:          "2025-06-01",
				Time:          "14:00",
				People:        5,
				Location:      "Test Location",
				CoffeeFlavors: []string{"french_toast"},
				MilkOptions:   []string{"whole"},
				Archived:      false,
			},
			expectedValid: false,
			fieldErrors: map[string]bool{
				"Name": true,
			},
		},
		{
			name: "Zero people",
			booking: models.Booking{
				Name:          "Test User",
				Email:         "test@example.com",
				Date:          "2025-06-01",
				Time:          "14:00",
				People:        0, // Should be at least 1
				Location:      "Test Location",
				CoffeeFlavors: []string{"french_toast"},
				MilkOptions:   []string{"whole"},
				Archived:      false,
			},
			expectedValid: false,
			fieldErrors: map[string]bool{
				"People": true,
			},
		},
		{
			name: "Empty coffee flavors",
			booking: models.Booking{
				Name:          "Test User",
				Email:         "test@example.com",
				Date:          "2025-06-01",
				Time:          "14:00",
				People:        5,
				Location:      "Test Location",
				CoffeeFlavors: []string{},
				MilkOptions:   []string{"whole"},
				Archived:      false,
			},
			expectedValid: false,
			fieldErrors: map[string]bool{
				"CoffeeFlavors": true,
			},
		},
		{
			name: "Empty milk options",
			booking: models.Booking{
				Name:          "Test User",
				Email:         "test@example.com",
				Date:          "2025-06-01",
				Time:          "14:00",
				People:        5,
				Location:      "Test Location",
				CoffeeFlavors: []string{"french_toast"},
				MilkOptions:   []string{},
				Archived:      false,
			},
			expectedValid: false,
			fieldErrors: map[string]bool{
				"MilkOptions": true,
			},
		},
		{
			name: "Contact validation - with email",
			booking: models.Booking{
				Name:          "Test User",
				Email:         "test@example.com",
				Phone:         "",
				Date:          "2025-06-01",
				Time:          "14:00",
				People:        5,
				Location:      "Test Location",
				CoffeeFlavors: []string{"french_toast"},
				MilkOptions:   []string{"whole"},
				Archived:      false,
			},
			expectedValid: true,
		},
		{
			name: "Contact validation - with phone",
			booking: models.Booking{
				Name:          "Test User",
				Email:         "",
				Phone:         "555-1234",
				Date:          "2025-06-01",
				Time:          "14:00",
				People:        5,
				Location:      "Test Location",
				CoffeeFlavors: []string{"french_toast"},
				MilkOptions:   []string{"whole"},
				Archived:      false,
			},
			expectedValid: true,
		},
		{
			name: "Contact validation - no email or phone",
			booking: models.Booking{
				Name:          "Test User",
				Email:         "",
				Phone:         "",
				Date:          "2025-06-01",
				Time:          "14:00",
				People:        5,
				Location:      "Test Location",
				CoffeeFlavors: []string{"french_toast"},
				MilkOptions:   []string{"whole"},
				Archived:      false,
			},
			expectedValid: false,
			fieldErrors: map[string]bool{
				"Email": true,
				"Phone": true,
			},
		},
	}

	validate := validator.New()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validate.Struct(tc.booking)

			// Check if validation result matches expected
			if (err == nil) != tc.expectedValid {
				t.Errorf("Expected validation to be %v, but got err: %v", tc.expectedValid, err)
			}

			// If expecting specific field errors, check those
			if !tc.expectedValid && tc.fieldErrors != nil {
				if err == nil {
					t.Fatal("Expected validation errors but got none")
				}

				validationErrors, ok := err.(validator.ValidationErrors)
				if !ok {
					t.Fatalf("Expected validator.ValidationErrors, got %T", err)
				}

				// Check each expected field error
				for field := range tc.fieldErrors {
					found := false
					for _, fieldErr := range validationErrors {
						if fieldErr.Field() == field {
							found = true
							break
						}
					}

					if !found {
						t.Errorf("Expected validation error for field %q, but none found", field)
					}
				}
			}
		})
	}
}

// Custom validation for contact information
func TestBookingContactValidation(t *testing.T) {
	booking := models.Booking{
		Name:          "Test User",
		Email:         "",
		Phone:         "",
		Date:          "2025-06-01",
		Time:          "14:00",
		People:        5,
		Location:      "Test Location",
		CoffeeFlavors: []string{"french_toast"},
		MilkOptions:   []string{"whole"},
		Archived:      false,
	}

	// Create validator with struct-level validation for contacts
	validate := validator.New()
	validate.RegisterStructValidation(func(sl validator.StructLevel) {
		b := sl.Current().Interface().(models.Booking)
		if b.Email == "" && b.Phone == "" {
			sl.ReportError(b.Email, "Email", "Email", "required_without_phone", "")
			sl.ReportError(b.Phone, "Phone", "Phone", "required_without_email", "")
		}
	}, models.Booking{})

	// Validate the booking with no contact info - should fail
	err := validate.Struct(booking)
	if err == nil {
		t.Error("Booking with no email or phone should be invalid")
	}

	// Add email and validate again - should pass
	booking.Email = "test@example.com"
	err = validate.Struct(booking)
	if err != nil {
		t.Errorf("Booking with email should be valid, got error: %v", err)
	}
}

// Booking Archive/Unarchive rules
func TestBookingArchiveRules(t *testing.T) {
	tests := []struct {
		name          string
		booking       models.Booking
		expectedValid bool
		fieldErrors   map[string]bool
	}{
		{
			name: "New booking not archived by default",
			booking: models.Booking{
				Name:          "Test User",
				Email:         "test@example.com",
				Date:          "2025-06-01", // Future date
				Time:          "14:00",
				People:        5,
				Location:      "Test Location",
				CoffeeFlavors: []string{"french_toast"},
				MilkOptions:   []string{"whole"},
				Archived:      false, // Not archived
			},
			expectedValid: true,
		},
		{
			name: "Archived past booking is valid",
			booking: models.Booking{
				Name:          "Test User",
				Email:         "test@example.com",
				Date:          "2024-01-01", // Past date
				Time:          "14:00",
				People:        5,
				Location:      "Test Location",
				CoffeeFlavors: []string{"french_toast"},
				MilkOptions:   []string{"whole"},
				Archived:      true, // Archived
			},
			expectedValid: true,
		},
		{
			name: "Unarchived past booking is valid",
			booking: models.Booking{
				Name:          "Test User",
				Email:         "test@example.com",
				Date:          "2024-01-01", // Past date
				Time:          "14:00",
				People:        5,
				Location:      "Test Location",
				CoffeeFlavors: []string{"french_toast"},
				MilkOptions:   []string{"whole"},
				Archived:      false, // Not archived
			},
			expectedValid: true,
		},
	}

	validate := validator.New()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validate.Struct(tc.booking)

			// Check if validation result matches expected
			if (err == nil) != tc.expectedValid {
				t.Errorf("Expected validation to be %v, but got err: %v", tc.expectedValid, err)
			}

			// If expecting specific field errors, check those
			if !tc.expectedValid && tc.fieldErrors != nil {
				if err == nil {
					t.Fatal("Expected validation errors but got none")
				}

				validationErrors, ok := err.(validator.ValidationErrors)
				if !ok {
					t.Fatalf("Expected validator.ValidationErrors, got %T", err)
				}

				// Check each expected field error
				for field := range tc.fieldErrors {
					found := false
					for _, fieldErr := range validationErrors {
						if fieldErr.Field() == field {
							found = true
							break
						}
					}

					if !found {
						t.Errorf("Expected validation error for field %q, but none found", field)
					}
				}
			}
		})
	}

}

// TestBookingArchiveBusinessRules tests more complex business rules for archiving
func TestBookingArchiveBusinessRules(t *testing.T) {
	// Current date for testing
	currentTime := time.Now()
	pastTime := currentTime.AddDate(0, -1, 0)  // 1 month ago
	futureTime := currentTime.AddDate(0, 1, 0) // 1 month from now

	// Format dates as strings like your API uses
	pastDateStr := pastTime.Format("2006-01-02")
	futureDateStr := futureTime.Format("2006-01-02")

	tests := []struct {
		name         string
		booking      models.Booking
		canArchive   bool
		canUnarchive bool
	}{
		{
			name: "Past booking can be archived",
			booking: models.Booking{
				Name:     "Past Booking",
				Email:    "past@example.com",
				Date:     pastDateStr,
				Archived: false,
			},
			canArchive:   true,
			canUnarchive: false, // Can't unarchive what's not archived
		},
		{
			name: "Archived booking can be unarchived",
			booking: models.Booking{
				Name:     "Archived Booking",
				Email:    "archived@example.com",
				Date:     pastDateStr,
				Archived: true,
			},
			canArchive:   false, // Already archived
			canUnarchive: true,
		},
		{
			name: "Future booking can be archived",
			booking: models.Booking{
				Name:     "Future Booking",
				Email:    "future@example.com",
				Date:     futureDateStr,
				Archived: false,
			},
			canArchive:   true, // Our current implementation allows this
			canUnarchive: false,
		},
		{
			name: "Already archived future booking cannot be archived again",
			booking: models.Booking{
				Name:     "Active Future Booking",
				Email:    "active@example.com",
				Date:     futureDateStr,
				Archived: true,
			},
			canArchive:   false, // Already archived
			canUnarchive: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Check CanArchiveBooking
			if models.CanArchiveBooking(&tc.booking) != tc.canArchive {
				t.Errorf("CanArchiveBooking returned %v, want %v",
					models.CanArchiveBooking(&tc.booking), tc.canArchive)
			}

			// Check CanUnarchiveBooking
			if models.CanUnarchiveBooking(&tc.booking) != tc.canUnarchive {
				t.Errorf("CanUnarchiveBooking returned %v, want %v",
					models.CanUnarchiveBooking(&tc.booking), tc.canUnarchive)
			}
		})
	}
}
