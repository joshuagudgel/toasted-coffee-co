package models

import (
	"time"
)

type ItemType string

const (
	CoffeeFlavor ItemType = "coffee_flavor"
	MilkOption   ItemType = "milk_option"
)

// MenuItem represents a menu item (coffee flavor or milk option)
type MenuItem struct {
	ID        int       `json:"id,omitempty"`
	Value     string    `json:"value" validate:"required"`
	Label     string    `json:"label" validate:"required"`
	Type      ItemType  `json:"type" validate:"required,oneof=coffee_flavor milk_option"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}
