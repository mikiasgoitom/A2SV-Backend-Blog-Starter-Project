package entity

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Token represents an authentication token (access or refresh)
type Token struct {
	ID        string    `bson:"id" json:"id"`
	UserID    string    `bson:"user_id" json:"user_id"`
	TokenType TokenType `bson:"token_type" json:"token_type"`
	TokenHash string    `bson:"token_hash" json:"-"`
	ExpiresAt time.Time `bson:"expires_at" json:"expires_at"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	Revoke    bool      `bson:"revoke" json:"revoked"`
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
