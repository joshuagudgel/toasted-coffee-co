package database

import (
	"context"
	"errors"
	"fmt"
	"log"
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
        INSERT INTO bookings (name, email, phone, date, time, people, location, notes, coffee_flavors, milk_options, package)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING id
    `, booking.Name, booking.Email, booking.Phone, parsedDate, booking.Time, booking.People, booking.Location,
		booking.Notes, booking.CoffeeFlavors, booking.MilkOptions, booking.Package).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetByID retrieves a booking by its ID
func (r *BookingRepository) GetByID(ctx context.Context, id int) (*models.Booking, error) {
	booking := &models.Booking{}

	var dateTime time.Time

	err := r.db.Pool.QueryRow(ctx, `
        SELECT id, name, email, phone, date, time, people, location, notes, coffee_flavors, milk_options, package, created_at 
        FROM bookings 
        WHERE id = $1
    `, id).Scan(
		&booking.ID, &booking.Name, &booking.Email, &booking.Phone, &dateTime, &booking.Time, &booking.People,
		&booking.Location, &booking.Notes, &booking.CoffeeFlavors, &booking.MilkOptions,
		&booking.Package, &booking.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Assign the date as a string in YYYY-MM-DD format
	booking.Date = dateTime.Format("2006-01-02")

	return booking, nil
}

// GetAll retrieves all bookings
func (r *BookingRepository) GetAll(ctx context.Context) ([]*models.Booking, error) {
	log.Println("Starting GetAll query...")

	query := `
        SELECT id, name, email, phone, date, time, people, location, notes, coffee_flavors, milk_options, package, created_at 
        FROM bookings
        ORDER BY date DESC
    `
	log.Println("Executing query:", query)

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		log.Printf("Database query error: %v", err)
		return nil, fmt.Errorf("database query error: %w", err)
	}
	defer rows.Close()

	// Initialize as empty slice rather than nil to ensure we return [] instead of null
	bookings := []*models.Booking{}
	rowNum := 0

	for rows.Next() {
		rowNum++
		booking := &models.Booking{}

		var dateTime time.Time // Temporary variable for date

		err := rows.Scan(
			&booking.ID, &booking.Name, &booking.Email, &booking.Phone, &dateTime, &booking.Time, &booking.People,
			&booking.Location, &booking.Notes, &booking.CoffeeFlavors, &booking.MilkOptions,
			&booking.Package, &booking.CreatedAt,
		)
		if err != nil {
			log.Printf("Error scanning row %d: %v", rowNum, err)
			return nil, fmt.Errorf("error scanning row %d: %w", rowNum, err)
		}

		// Assign the date
		booking.Date = dateTime.Format("2006-01-02")

		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error after scanning rows: %v", err)
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	log.Printf("Successfully retrieved %d bookings", len(bookings))
	return bookings, nil
}

// Delete removes a booking from the database
func (r *BookingRepository) Delete(ctx context.Context, id int) error {
	// Execute the delete query
	commandTag, err := r.db.Pool.Exec(ctx, `
        DELETE FROM bookings 
        WHERE id = $1
    `, id)

	if err != nil {
		return err
	}

	// Check if any rows were affected
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("booking not found")
	}

	return nil
}

// Update modifies an existing booking
func (r *BookingRepository) Update(ctx context.Context, id int, booking *models.Booking) error {
	// Parse date string to time.Time
	parsedDate, err := time.Parse("2006-01-02", booking.Date)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}

	commandTag, err := r.db.Pool.Exec(ctx, `
        UPDATE bookings 
        SET name = $1, email = $2, phone = $3, date = $4, time = $5, 
            people = $6, location = $7, notes = $8, coffee_flavors = $9, 
            milk_options = $10, package = $11
        WHERE id = $12
    `, booking.Name, booking.Email, booking.Phone, parsedDate, booking.Time,
		booking.People, booking.Location, booking.Notes, booking.CoffeeFlavors,
		booking.MilkOptions, booking.Package, id)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("booking not found")
	}

	return nil
}
