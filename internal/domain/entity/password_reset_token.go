package entity

import (
	"time"
)

// PasswordResetToken represents a password reset token
type PasswordResetToken struct {
	ID        string    `json:"id" bson:"_id"`
	UserID    string    `json:"user_id" bson:"user_id"`
	TokenHash string    `json:"-" bson:"token_hash"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`
	Used      bool      `json:"used" bson:"used"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
