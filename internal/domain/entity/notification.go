package entity

import (
	"time"
)

// Notification represents a notification sent to a user
type Notification struct {
	ID              string           `json:"id" db:"id"`
	RecipientUserID string           `json:"recipient_user_id" db:"recipient_user_id"`
	SenderUserID    *string          `json:"sender_user_id" db:"sender_user_id"`
	Type            NotificationType `json:"type" db:"type"`
	Message         string           `json:"message" db:"message"`
	RelatedEntityID *string          `json:"related_entity_id" db:"related_entity_id"`
	IsRead          bool             `json:"is_read" db:"is_read"`
	CreatedAt       time.Time        `json:"created_at" db:"created_at"`
}

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationTypeNewComment        NotificationType = "NEW_COMMENT"
	NotificationTypePostLiked         NotificationType = "POST_LIKED"
	NotificationTypePasswordReset     NotificationType = "PASSWORD_RESET"
	NotificationTypeEmailVerification NotificationType = "EMAIL_VERIFICATION"
	NotificationTypeCommentLiked      NotificationType = "COMMENT_LIKED"
	NotificationTypePackageExpired    NotificationType = "PACKAGE_EXPIRED"
)
