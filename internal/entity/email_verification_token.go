package entity

import (
	"time"

	"github.com/google/uuid"
)

// EmailVerificationToken represents an email verification token
type EmailVerificationToken struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	TokenHash string    `json:"-" db:"token_hash"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	Used      bool      `json:"used" db:"used"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// // NewEmailVerificationToken creates a new email verification token
// func NewEmailVerificationToken(userID uuid.UUID, tokenHash string, expiresAt time.Time) *EmailVerificationToken {
// 	return &EmailVerificationToken{
// 		ID:        uuid.New(),
// 		UserID:    userID,
// 		TokenHash: tokenHash,
// 		ExpiresAt: expiresAt,
// 		Used:      false,
// 		CreatedAt: time.Now(),
// 	}
// }

// // MarkAsUsed marks the token as used
// func (e *EmailVerificationToken) MarkAsUsed() {
// 	e.Used = true
// }

// // IsExpired checks if the token has expired
// func (e *EmailVerificationToken) IsExpired() bool {
// 	return time.Now().After(e.ExpiresAt)
// }

// // IsValid checks if the token is valid (not expired and not used)
// func (e *EmailVerificationToken) IsValid() bool {
// 	return !e.IsExpired() && !e.Used
// } 