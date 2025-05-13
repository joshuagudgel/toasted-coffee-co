package models

import (
	"time"
)

// Booking represents a coffee booking
type Booking struct {
	ID            int       `json:"id,omitempty"`
	Name          string    `json:"name" validate:"required"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Date          string    `json:"date" validate:"required"`
	Time          string    `json:"time" validate:"required"`
	People        int       `json:"people" validate:"required,min=1"`
	Location      string    `json:"location" validate:"required"`
	Notes         string    `json:"notes"`
	CoffeeFlavors []string  `json:"coffeeFlavors" validate:"required,min=1"`
	MilkOptions   []string  `json:"milkOptions" validate:"required,min=1"`
	Package       string    `json:"package"`
	CreatedAt     time.Time `json:"createdAt,omitempty"`
}
