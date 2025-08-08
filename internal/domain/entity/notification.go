package entity

import (
	"time"
)

// Notification represents a notification sent to a user
type Notification struct {
	ID              string           `json:"id" bson:"_id"`
	RecipientUserID string           `json:"recipient_user_id" bson:"recipient_user_id"`
	SenderUserID    *string          `json:"sender_user_id" bson:"sender_user_id"`
	Type            NotificationType `json:"type" bson:"type"`
	Message         string           `json:"message" bson:"message"`
	RelatedEntityID *string          `json:"related_entity_id" bson:"related_entity_id"`
	IsRead          bool             `json:"is_read" bson:"is_read"`
	CreatedAt       time.Time        `json:"created_at" bson:"created_at"`
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
