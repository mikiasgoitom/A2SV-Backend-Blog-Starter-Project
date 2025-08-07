package entity

import (
	"time"
)

// AIQuery represents an AI interaction/query
type AIQuery struct {
	ID        string    `json:"id" bson:"_id"`
	UserID    *string   `json:"user_id" bson:"user_id"`
	Prompt    string    `json:"prompt" bson:"prompt"`
	Response  string    `json:"response" bson:"response"`
	ModelUsed string    `json:"model_used" bson:"model_used"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}
