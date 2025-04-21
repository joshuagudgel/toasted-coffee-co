package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joshuagudgel/toasted-coffee/backend/internal/models"
)

// BookingRepository handles database operations for bookings
type BookingRepository struct {
	db *DB
}

// NewBookingRepository creates a new booking repository
func NewBookingRepository(db *DB) *BookingRepository {
	return &BookingRepository{db: db}
}

// Create inserts a new booking into the database
func (r *BookingRepository) Create(ctx context.Context, booking *models.Booking) (int, error) {
	// Parse date string to time.Time if needed
	parsedDate, err := time.Parse("2006-01-02", booking.Date)
	if err != nil {
		return 0, fmt.Errorf("invalid date format: %w", err)
	}

	var id int
	err = r.db.Pool.QueryRow(ctx, `
        INSERT INTO bookings (name, date, time, people, location, notes, coffee_flavors, milk_options, package)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id
    `, booking.Name, parsedDate, booking.Time, booking.People, booking.Location,
		booking.Notes, booking.CoffeeFlavors, booking.MilkOptions, booking.Package).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetByID retrieves a booking by its ID
func (r *BookingRepository) GetByID(ctx context.Context, id int) (*models.Booking, error) {
	booking := &models.Booking{}

	err := r.db.Pool.QueryRow(ctx, `
        SELECT id, name, date, time, people, location, notes, coffee_flavors, milk_options, package, created_at 
        FROM bookings 
        WHERE id = $1
    `, id).Scan(
		&booking.ID, &booking.Name, &booking.Date, &booking.Time, &booking.People,
		&booking.Location, &booking.Notes, &booking.CoffeeFlavors, &booking.MilkOptions,
		&booking.Package, &booking.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return booking, nil
}

// GetAll retrieves all bookings
func (r *BookingRepository) GetAll(ctx context.Context) ([]*models.Booking, error) {
	rows, err := r.db.Pool.Query(ctx, `
        SELECT id, name, date, time, people, location, notes, coffee_flavors, milk_options, package, created_at 
        FROM bookings
        ORDER BY date DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*models.Booking
	for rows.Next() {
		booking := &models.Booking{}
		if err := rows.Scan(
			&booking.ID, &booking.Name, &booking.Date, &booking.Time, &booking.People,
			&booking.Location, &booking.Notes, &booking.CoffeeFlavors, &booking.MilkOptions,
			&booking.Package, &booking.CreatedAt,
		); err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}
