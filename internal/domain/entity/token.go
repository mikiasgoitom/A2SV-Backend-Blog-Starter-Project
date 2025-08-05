package entity

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Token represents an authentication token (access or refresh)
type Token struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	TokenType TokenType `json:"token_type" db:"token_type"`
	TokenHash string    `json:"-" db:"token_hash"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Revoke    bool      `json:"revoked" db:"revoked"`
}

// TokenType represents the type of token
type TokenType string

const (
	TokenTypeAccess            TokenType = "access"
	TokenTypeRefresh           TokenType = "refresh"
	TokenTypePasswordReset     TokenType = "password_reset"
	TokenTypeEmailVerification TokenType = "email_verification"
)

func isValidTokenType(tokType string) bool {

	switch TokenType(tokType) {
	case TokenTypeAccess, TokenTypeRefresh, TokenTypePasswordReset:
		return true
	default:
		return false
	}
}

func SetTokenType(tokType string) (TokenType, error) {
	if isValidTokenType(tokType) {
		return TokenType(tokType), nil
	} else {
		return "", fmt.Errorf("invalid token type: %s", tokType)
	}
}

// Claims represents the JWT claims for authentication and authorization.
type Claims struct {
	UserID string   `json:"user_id"`
	Role   UserRole `json:"role"`
	jwt.RegisteredClaims
}
