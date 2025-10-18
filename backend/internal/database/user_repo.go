package database

import (
	"context"

	"github.com/joshuagudgel/toasted-coffee/backend/internal/models"
)

type UserRepository struct {
	db *DB
}

// NewUserRepository creates a new menu repository
func NewUserRepository(db *DB) UserRepositoryInterface {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	user := &models.User{}

	err := r.db.Pool.QueryRow(ctx, `
        SELECT id, username, password, role FROM users WHERE id = $1
    `, id).Scan(&user.ID, &user.Username, &user.Password, &user.Role)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}

	err := r.db.Pool.QueryRow(ctx, `
        SELECT id, username, password, role FROM users WHERE username = $1
    `, username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)

	if err != nil {
		return nil, err
	}

	return user, nil
}
