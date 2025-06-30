package database

import (
	"context"
	"time"

	"github.com/joshuagudgel/toasted-coffee/backend/internal/models"
)

// PackageRepository handles database operations for service packages
type PackageRepository interface {
	GetAll(ctx context.Context, includeInactive bool) ([]models.Package, error)
	GetByID(ctx context.Context, id int) (*models.Package, error)
	Create(ctx context.Context, pkg *models.PackageInput) (int, error)
	Update(ctx context.Context, id int, pkg *models.PackageInput) error
	Delete(ctx context.Context, id int) error
}

type packageRepository struct {
	db *DB
}

// NewPackageRepository creates a new package repository
func NewPackageRepository(db *DB) PackageRepository {
	return &packageRepository{db: db}
}

// GetAll retrieves all packages
func (r *packageRepository) GetAll(ctx context.Context, includeInactive bool) ([]models.Package, error) {
	query := `
        SELECT p.id, p.name, p.price, p.description, p.display_order, p.active, p.created_at, p.updated_at
        FROM packages p
    `

	if !includeInactive {
		query += " WHERE p.active = true"
	}

	query += " ORDER BY p.display_order, p.name"

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	packages := []models.Package{}
	for rows.Next() {
		var pkg models.Package
		if err := rows.Scan(
			&pkg.ID,
			&pkg.Name,
			&pkg.Price,
			&pkg.Description,
			&pkg.DisplayOrder,
			&pkg.Active,
			&pkg.CreatedAt,
			&pkg.UpdatedAt,
		); err != nil {
			return nil, err
		}

		// Get points for this package
		pointsQuery := `
            SELECT point FROM package_points 
            WHERE package_id = $1 
            ORDER BY display_order
        `
		pointRows, err := r.db.Pool.Query(ctx, pointsQuery, pkg.ID)
		if err != nil {
			return nil, err
		}
		defer pointRows.Close()

		points := []string{}
		for pointRows.Next() {
			var point string
			if err := pointRows.Scan(&point); err != nil {
				return nil, err
			}
			points = append(points, point)
		}
		pkg.Points = points
		packages = append(packages, pkg)
	}

	return packages, nil
}

// GetByID retrieves a package by ID
func (r *packageRepository) GetByID(ctx context.Context, id int) (*models.Package, error) {
	query := `
        SELECT id, name, price, description, display_order, active, created_at, updated_at
        FROM packages
        WHERE id = $1
    `

	var pkg models.Package
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&pkg.ID,
		&pkg.Name,
		&pkg.Price,
		&pkg.Description,
		&pkg.DisplayOrder,
		&pkg.Active,
		&pkg.CreatedAt,
		&pkg.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Get points for this package
	pointsQuery := `
        SELECT point FROM package_points 
        WHERE package_id = $1 
        ORDER BY display_order
    `
	rows, err := r.db.Pool.Query(ctx, pointsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	points := []string{}
	for rows.Next() {
		var point string
		if err := rows.Scan(&point); err != nil {
			return nil, err
		}
		points = append(points, point)
	}
	pkg.Points = points

	return &pkg, nil
}

// Create adds a new package
func (r *packageRepository) Create(ctx context.Context, input *models.PackageInput) (int, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var packageID int
	err = tx.QueryRow(ctx, `
        INSERT INTO packages (name, price, description, display_order, active, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `, input.Name, input.Price, input.Description, input.DisplayOrder, input.Active, time.Now()).Scan(&packageID)
	if err != nil {
		return 0, err
	}

	// Insert points
	for i, point := range input.Points {
		_, err = tx.Exec(ctx, `
            INSERT INTO package_points (package_id, point, display_order)
            VALUES ($1, $2, $3)
        `, packageID, point, i)
		if err != nil {
			return 0, err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}

	return packageID, nil
}

// Update modifies an existing package
func (r *packageRepository) Update(ctx context.Context, id int, input *models.PackageInput) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Update package
	_, err = tx.Exec(ctx, `
        UPDATE packages
        SET name = $1, price = $2, description = $3, display_order = $4, active = $5, updated_at = $6
        WHERE id = $7
    `, input.Name, input.Price, input.Description, input.DisplayOrder, input.Active, time.Now(), id)
	if err != nil {
		return err
	}

	// Delete existing points
	_, err = tx.Exec(ctx, `DELETE FROM package_points WHERE package_id = $1`, id)
	if err != nil {
		return err
	}

	// Insert new points
	for i, point := range input.Points {
		_, err = tx.Exec(ctx, `
            INSERT INTO package_points (package_id, point, display_order)
            VALUES ($1, $2, $3)
        `, id, point, i)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// Delete removes a package
func (r *packageRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.Pool.Exec(ctx, `DELETE FROM packages WHERE id = $1`, id)
	return err
}
