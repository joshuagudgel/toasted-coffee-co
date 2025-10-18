package database

import (
	"context"

	"github.com/joshuagudgel/toasted-coffee/backend/internal/models"
)

type Repositories struct {
	Booking BookingRepositoryInterface
	User    UserRepositoryInterface
	Menu    MenuRepositoryInterface
	Package PackageRepositoryInterface
}

// BookingRepositoryInterface defines the methods for booking operations
type BookingRepositoryInterface interface {
	Create(ctx context.Context, booking *models.Booking) (int, error)
	GetByID(ctx context.Context, id int) (*models.Booking, error)
	GetAll(ctx context.Context, includeArchived bool) ([]*models.Booking, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, booking *models.Booking) error
	Archive(ctx context.Context, id int) error
	Unarchive(ctx context.Context, id int) error
}

// UserRepositoryInterface defines the methods for user operations
type UserRepositoryInterface interface {
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
}

// MenuRespositoryInterface defines the methods for menu operations
type MenuRepositoryInterface interface {
	GetAll(ctx context.Context) ([]models.MenuItem, error)
	GetByType(ctx context.Context, itemType models.ItemType) ([]models.MenuItem, error)
	Create(ctx context.Context, item *models.MenuItem) (int, error)
	Update(ctx context.Context, id int, item *models.MenuItem) error
	Delete(ctx context.Context, id int) error
}

// PackageRepositoryInterface defines the methods for package operations
type PackageRepositoryInterface interface {
	GetAll(ctx context.Context, includeInactive bool) ([]models.Package, error)
	GetByID(ctx context.Context, id int) (*models.Package, error)
	Create(ctx context.Context, pkg *models.PackageInput) (int, error)
	Update(ctx context.Context, id int, pkg *models.PackageInput) error
	Delete(ctx context.Context, id int) error
}

// NewRepositories creates all repositories
func NewRepositories(db *DB) *Repositories {
	return &Repositories{
		Booking: NewBookingRepository(db),
		User:    NewUserRepository(db),
		Menu:    NewMenuRepository(db),
		Package: NewPackageRepository(db),
	}
}
