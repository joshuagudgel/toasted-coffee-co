package database

import (
	"context"

	"github.com/joshuagudgel/toasted-coffee/backend/internal/models"
)

// MenuRepository handles database operations for menu items
type MenuRepository interface {
	GetAll(ctx context.Context) ([]models.MenuItem, error)
	GetByType(ctx context.Context, itemType models.ItemType) ([]models.MenuItem, error)
	Create(ctx context.Context, item *models.MenuItem) (int, error)
	Update(ctx context.Context, id int, item *models.MenuItem) error
	Delete(ctx context.Context, id int) error
}

type menuRepository struct {
	db *DB
}

// NewMenuRepository creates a new menu repository
func NewMenuRepository(db *DB) MenuRepository {
	return &menuRepository{
		db: db,
	}
}

// Implementation of repository methods
func (r *menuRepository) GetAll(ctx context.Context) ([]models.MenuItem, error) {
	rows, err := r.db.Pool.Query(ctx, `
        SELECT id, value, label, type, active, created_at, updated_at
        FROM menu_items
        ORDER BY type, label
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		var itemType string
		if err := rows.Scan(
			&item.ID, &item.Value, &item.Label, &itemType, &item.Active,
			&item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		item.Type = models.ItemType(itemType)
		items = append(items, item)
	}

	return items, nil
}

// GetByType retrieves menu items of a specific type
func (r *menuRepository) GetByType(ctx context.Context, itemType models.ItemType) ([]models.MenuItem, error) {
    rows, err := r.db.Pool.Query(ctx, `
        SELECT id, value, label, type, active, created_at, updated_at
        FROM menu_items
        WHERE type = $1
        ORDER BY label
    `, string(itemType))
    
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []models.MenuItem
    for rows.Next() {
        var item models.MenuItem
        var itemType string
        if err := rows.Scan(
            &item.ID, &item.Value, &item.Label, &itemType, &item.Active,
            &item.CreatedAt, &item.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        item.Type = models.ItemType(itemType)
        items = append(items, item)
    }

    return items, nil
}

// Create adds a new menu item
func (r *menuRepository) Create(ctx context.Context, item *models.MenuItem) (int, error) {
	var id int
	err := r.db.Pool.QueryRow(ctx, `
        INSERT INTO menu_items (value, label, type, active)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `, item.Value, item.Label, item.Type, item.Active).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update modifies an existing menu item
func (r *menuRepository) Update(ctx context.Context, id int, item *models.MenuItem) error {
	_, err := r.db.Pool.Exec(ctx, `
        UPDATE menu_items
        SET value = $1, label = $2, type = $3, active = $4, updated_at = CURRENT_TIMESTAMP
        WHERE id = $5
    `, item.Value, item.Label, item.Type, item.Active, id)

	return err
}

// Delete removes a menu item
func (r *menuRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.Pool.Exec(ctx, `
        DELETE FROM menu_items
        WHERE id = $1
    `, id)

	return err
}
