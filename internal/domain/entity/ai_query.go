package entity

import (
	"time"
	"github.com/google/uuid"
)

// AIQuery represents an AI interaction/query
type AIQuery struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    *uuid.UUID `json:"user_id" db:"user_id"`
	Prompt    string    `json:"prompt" db:"prompt"`
	Response  string    `json:"response" db:"response"`
	ModelUsed string    `json:"model_used" db:"model_used"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}
