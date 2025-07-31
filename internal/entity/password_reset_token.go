package entity

import (
	"time"

	"github.com/google/uuid"
)

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	TokenHash string    `json:"-" db:"token_hash"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	Used      bool      `json:"used" db:"used"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// // NewPasswordResetToken creates a new password reset token
// func NewPasswordResetToken(userID uuid.UUID, tokenHash string, expiresAt time.Time) *PasswordResetToken {
// 	return &PasswordResetToken{
// 		ID:        uuid.New(),
// 		UserID:    userID,
// 		TokenHash: tokenHash,
// 		ExpiresAt: expiresAt,
// 		Used:      false,
// 		CreatedAt: time.Now(),
// 	}
// }

// // MarkAsUsed marks the token as used
// func (p *PasswordResetToken) MarkAsUsed() {
// 	p.Used = true
// }

// // IsExpired checks if the token has expired
// func (p *PasswordResetToken) IsExpired() bool {
// 	return time.Now().After(p.ExpiresAt)
// }

// // IsValid checks if the token is valid (not expired and not used)
// func (p *PasswordResetToken) IsValid() bool {
// 	return !p.IsExpired() && !p.Used
// } 