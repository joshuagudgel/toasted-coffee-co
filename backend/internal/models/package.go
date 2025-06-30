package models

import "time"

// Package represents a service package offered by the company
type Package struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Price        string    `json:"price"`
	Description  string    `json:"description"`
	Points       []string  `json:"points"`
	DisplayOrder int       `json:"displayOrder"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// PackageInput is used for creating or updating packages
type PackageInput struct {
	Name         string   `json:"name"`
	Price        string   `json:"price"`
	Description  string   `json:"description"`
	Points       []string `json:"points"`
	DisplayOrder int      `json:"displayOrder"`
	Active       bool     `json:"active"`
}
