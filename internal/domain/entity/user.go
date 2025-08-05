package entity

import (
	"time"

	"github.com/google/uuid"
)

// User represents a registered user in the system
type User struct {
	ID           uuid.UUID `db:"id"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Role         UserRole  `db:"role"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	FirstName    *string   `db:"first_name"`
	LastName     *string   `db:"last_name"`
	AvatarURL    *string   `db:"avatar_url"`
}

// UserRole represents the role of a user in the system
type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)

func DefaultRole() UserRole {
	return UserRoleUser
}
