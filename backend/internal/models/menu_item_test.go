package models_test

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/models"
)

func TestMenuItemValidation(t *testing.T) {
	tests := []struct {
		name          string
		menuItem      models.MenuItem
		expectedValid bool
		fieldErrors   map[string]bool // Fields expected to have validation errors
	}{
		{
			name: "Valid coffee flavor",
			menuItem: models.MenuItem{
				Value:  "french_toast",
				Label:  "French Toast",
				Type:   models.CoffeeFlavor,
				Active: true,
			},
			expectedValid: true,
		},
		{
			name: "Valid milk option",
			menuItem: models.MenuItem{
				Value:  "whole",
				Label:  "Whole Milk",
				Type:   models.MilkOption,
				Active: true,
			},
			expectedValid: true,
		},
		{
			name: "Missing value",
			menuItem: models.MenuItem{
				Value:  "", // missing
				Label:  "French Toast",
				Type:   models.CoffeeFlavor,
				Active: true,
			},
			expectedValid: false,
			fieldErrors: map[string]bool{
				"Value": true,
			},
		},
		{
			name: "Missing label",
			menuItem: models.MenuItem{
				Value:  "french_toast",
				Label:  "", // missing
				Type:   models.CoffeeFlavor,
				Active: true,
			},
			expectedValid: false,
			fieldErrors: map[string]bool{
				"Label": true,
			},
		},
		{
			name: "Invalid type",
			menuItem: models.MenuItem{
				Value:  "french_toast",
				Label:  "French Toast",
				Type:   "invalid_type", // invalid
				Active: true,
			},
			expectedValid: false,
			fieldErrors: map[string]bool{
				"Type": true,
			},
		},
	}

	validate := validator.New()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validate.Struct(tc.menuItem)

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

func TestItemTypeConstants(t *testing.T) {
	// Test that the constants have the expected values
	if models.CoffeeFlavor != "coffee_flavor" {
		t.Errorf("Expected CoffeeFlavor to be 'coffee_flavor', got '%s'", models.CoffeeFlavor)
	}

	if models.MilkOption != "milk_option" {
		t.Errorf("Expected MilkOption to be 'milk_option', got '%s'", models.MilkOption)
	}
}
