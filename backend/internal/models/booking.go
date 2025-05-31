package models

import (
	"time"
)

// Booking represents a coffee booking
type Booking struct {
	ID            int       `json:"id,omitempty"`
	Name          string    `json:"name" validate:"required"`
	Email         string    `json:"email" validate:"required_without=Phone,omitempty"`
	Phone         string    `json:"phone" validate:"required_without=Email,omitempty"`
	Date          string    `json:"date" validate:"required"`
	Time          string    `json:"time" validate:"required"`
	People        int       `json:"people" validate:"required,min=1"`
	Location      string    `json:"location" validate:"required"`
	Notes         string    `json:"notes"`
	CoffeeFlavors []string  `json:"coffeeFlavors" validate:"required,min=1"`
	MilkOptions   []string  `json:"milkOptions" validate:"required,min=1"`
	Package       string    `json:"package"`
	CreatedAt     time.Time `json:"createdAt,omitempty"`
	Archived      bool      `json:"archived"`
}

// CanArchiveBooking determines if a booking can be archived
func CanArchiveBooking(booking *Booking) bool {
	// Parse the booking date
	bookingDate, err := time.Parse("2006-01-02", booking.Date)
	if (err != nil) {
		return false
	}
	
	// Current date
	currentDate := time.Now()
	
	// Rule 1: Already archived bookings can't be archived again
	if booking.Archived {
		return false
	}
	
	// Rule 2: Past bookings can always be archived
	if bookingDate.Before(currentDate) {
		return true
	}
	
	// Rule 3: Future bookings can be archived if they have status "canceled"
	// Note: This is commented out as you may not have this field yet
	// if booking.Status == "canceled" {
	//     return true
	// }
	
	// For now, allow archiving any booking
	return true
}

// CanUnarchiveBooking determines if a booking can be unarchived
func CanUnarchiveBooking(booking *Booking) bool {
	// Can only unarchive bookings that are currently archived
	return booking.Archived
}
