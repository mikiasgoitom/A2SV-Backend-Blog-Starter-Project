package entity

import (
	"time"

	"github.com/google/uuid"
)

// Token represents an authentication token (access or refresh)
type Token struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	TokenType TokenType `json:"token_type" db:"token_type"`
	TokenHash string    `json:"-" db:"token_hash"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Revoked   bool      `json:"revoked" db:"revoked"`
}

// TokenType represents the type of token
type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)