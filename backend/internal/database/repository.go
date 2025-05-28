package database

import (
	"context"

	"github.com/joshuagudgel/toasted-coffee/backend/internal/models"
)

// BookingRepositoryInterface defines the methods for booking operations
type BookingRepositoryInterface interface {
	Create(ctx context.Context, booking *models.Booking) (int, error)
	GetByID(ctx context.Context, id int) (*models.Booking, error)
	GetAll(ctx context.Context) ([]*models.Booking, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, booking *models.Booking) error
}
