package models_test

import (
	"testing"

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
			},
			expectedValid: false,
			fieldErrors: map[string]bool{
				"Email": true,
				"Phone": true,
			},
		},
	}

	// You'll need to implement a validate function for your models
	// or use a validation library like go-validator
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
	}

	// This is the custom validation logic that your handler uses
	if booking.Email == "" && booking.Phone == "" {
		t.Error("Booking with no email or phone should be invalid")
	} else {
		t.Error("Test case should fail but passed")
	}
}
