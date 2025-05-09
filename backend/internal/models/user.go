package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username" validate:"required"`
	Password string `json:"-" validate:"required"`    // Never expose in JSON
	Role     string `json:"role" validate:"required"` // "admin" for full access
}
