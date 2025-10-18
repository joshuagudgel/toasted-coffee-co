package database

import (
	"context"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type Seeder struct {
	db *DB
}

func NewSeeder(db *DB) *Seeder {
	return &Seeder{db: db}
}

func (s *Seeder) SeedAdminUser() error {
	log.Println("Setting up admin user...")

	var count int
	err := s.db.Pool.QueryRow(context.Background(), `
        SELECT COUNT(*) FROM users WHERE username = $1
    `, "admin").Scan(&count)

	if err != nil {
		log.Printf("Warning: Failed to check for admin user: %v", err)
		return err
	}

	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		_, err = s.db.Pool.Exec(context.Background(), `
            INSERT INTO users (username, password, role) VALUES ($1, $2, $3)
        `, "admin", string(hashedPassword), "admin")

		if err != nil {
			return err
		}

		log.Println("Admin user created successfully")
	} else {
		log.Println("Admin user already exists, skipping creation")
	}

	return nil
}
