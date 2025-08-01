package entity

import "time"

/*
@MG790 - implement token interface and entity, implement the token_repository.go, write the auth middleware
*/

type Token struct {
	ID        string
	UserID    string
	TokenHash string
	CreatedAt time.Time
	ExpiresAt time.Time
	Revoke    bool
}
