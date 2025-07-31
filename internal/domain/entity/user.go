package entity

import (
	"time"

	"github.com/google/uuid"
)

// User represents a registered user in the system
type User struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Username      string    `json:"username" db:"username"`
	Email         string    `json:"email" db:"email"`
	PasswordHash  string    `json:"-" db:"password_hash"`
	Role          UserRole  `json:"role" db:"role"`
	PackageID     *uuid.UUID `json:"package_id" db:"package_id"`
	PackageExpiry *time.Time `json:"package_expiry" db:"package_expiry"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	FirstName     *string   `json:"first_name" db:"first_name"`
	LastName      *string   `json:"last_name" db:"last_name"`
	AvatarURL     *string   `json:"avatar_url" db:"avatar_url"`
}

// UserRole represents the role of a user in the system
type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)