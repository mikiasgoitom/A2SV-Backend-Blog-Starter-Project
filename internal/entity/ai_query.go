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

// // NewAIQuery creates a new AI query instance
// func NewAIQuery(prompt, response, modelUsed string, userID *uuid.UUID) *AIQuery {
// 	return &AIQuery{
// 		ID:        uuid.New(),
// 		UserID:    userID,
// 		Prompt:    prompt,
// 		Response:  response,
// 		ModelUsed: modelUsed,
// 		Timestamp: time.Now(),
// 	}
// }

// // IsAnonymous checks if the query was made by an anonymous user
// func (a *AIQuery) IsAnonymous() bool {
// 	return a.UserID == nil
// }

// // GetQueryLength returns the length of the prompt
// func (a *AIQuery) GetQueryLength() int {
// 	return len(a.Prompt)
// }

// // GetResponseLength returns the length of the response
// func (a *AIQuery) GetResponseLength() int {
// 	return len(a.Response)
// } 