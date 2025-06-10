package handlers_test

import (
	"context"

	"github.com/joshuagudgel/toasted-coffee/backend/internal/models"
)

// MockMenuRepository implements MenuRepository interface for testing
type MockMenuRepository struct {
	// GetAll
	GetAllFunc   func(context.Context) ([]models.MenuItem, error)
	GetAllCalled bool

	// GetByType
	GetByTypeFunc   func(context.Context, models.ItemType) ([]models.MenuItem, error)
	GetByTypeCalled bool
	GetByTypeArg    models.ItemType

	// Create
	CreateFunc   func(context.Context, *models.MenuItem) (int, error)
	CreateCalled bool
	CreateItem   *models.MenuItem

	// Update
	UpdateFunc   func(context.Context, int, *models.MenuItem) error
	UpdateCalled bool
	UpdateID     int
	UpdateItem   *models.MenuItem

	// Delete
	DeleteFunc   func(context.Context, int) error
	DeleteCalled bool
	DeleteArg    int
}

// Implement interface methods
func (m *MockMenuRepository) GetAll(ctx context.Context) ([]models.MenuItem, error) {
	m.GetAllCalled = true
	return m.GetAllFunc(ctx)
}

func (m *MockMenuRepository) GetByType(ctx context.Context, itemType models.ItemType) ([]models.MenuItem, error) {
	m.GetByTypeCalled = true
	m.GetByTypeArg = itemType
	return m.GetByTypeFunc(ctx, itemType)
}

func (m *MockMenuRepository) Create(ctx context.Context, item *models.MenuItem) (int, error) {
	m.CreateCalled = true
	m.CreateItem = item
	return m.CreateFunc(ctx, item)
}

func (m *MockMenuRepository) Update(ctx context.Context, id int, item *models.MenuItem) error {
	m.UpdateCalled = true
	m.UpdateID = id
	m.UpdateItem = item
	return m.UpdateFunc(ctx, id, item)
}

func (m *MockMenuRepository) Delete(ctx context.Context, id int) error {
	m.DeleteCalled = true
	m.DeleteArg = id
	return m.DeleteFunc(ctx, id)
}
