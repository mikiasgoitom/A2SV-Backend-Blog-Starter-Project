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

// // NewToken creates a new token instance
// func NewToken(userID uuid.UUID, tokenType TokenType, tokenHash string, expiresAt time.Time) *Token {
// 	return &Token{
// 		ID:        uuid.New(),
// 		UserID:    userID,
// 		TokenType: tokenType,
// 		TokenHash: tokenHash,
// 		ExpiresAt: expiresAt,
// 		CreatedAt: time.Now(),
// 		Revoked:   false,
// 	}
// }

// // NewAccessToken creates a new access token
// func NewAccessToken(userID uuid.UUID, tokenHash string, expiresAt time.Time) *Token {
// 	return NewToken(userID, TokenTypeAccess, tokenHash, expiresAt)
// }

// // NewRefreshToken creates a new refresh token
// func NewRefreshToken(userID uuid.UUID, tokenHash string, expiresAt time.Time) *Token {
// 	return NewToken(userID, TokenTypeRefresh, tokenHash, expiresAt)
// }

// // Revoke marks the token as revoked
// func (t *Token) Revoke() {
// 	t.Revoked = true
// }

// // IsExpired checks if the token has expired
// func (t *Token) IsExpired() bool {
// 	return time.Now().After(t.ExpiresAt)
// }

// // IsValid checks if the token is valid (not expired and not revoked)
// func (t *Token) IsValid() bool {
// 	return !t.IsExpired() && !t.Revoked
// }

// // IsAccessToken checks if this is an access token
// func (t *Token) IsAccessToken() bool {
// 	return t.TokenType == TokenTypeAccess
// }

// // // IsRefreshToken checks if this is a refresh token
// // func (t *Token) IsRefreshToken() bool {