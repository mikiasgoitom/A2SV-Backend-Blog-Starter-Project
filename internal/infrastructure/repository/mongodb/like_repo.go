package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// LikeRepository represents the MongoDB implementation of the ILikeRepository interface.
type LikeRepository struct {
	collection *mongo.Collection
}

// NewLikeRepository creates and returns a new LikeRepository instance.
func NewLikeRepository(db *mongo.Database) *LikeRepository {
	return &LikeRepository{
		collection: db.Collection("likes"),
	}
}

func (r *LikeRepository) CreateReaction(ctx context.Context, like *entity.Like) error {
	// Filter to find an existing reaction by this user on this target.
	filter := bson.M{
		"user_id":     like.UserID,
		"target_id":   like.TargetID,
		"target_type": like.TargetType,
	}

	// Fields to set/update on the document.
	updateFields := bson.M{
		"type":       like.Type,  // Set the new type ("like" or "dislike")
		"is_deleted": false,      // Ensure it's marked as not deleted when created/updated
		"updated_at": time.Now(), // Update timestamp on any change
	}

	// Fields to set ONLY on initial insert (when upsert: true creates a new document)
	setOnInsertFields := bson.M{
		"id":         uuid.New(), // Generate a new ID only if inserting a new document
		"created_at": time.Now(), // Using camelCase to match entity db tags
	}

	updateDoc := bson.M{
		"$set":         updateFields,
		"$setOnInsert": setOnInsertFields,
	}

	opts := options.Update().SetUpsert(true)

	res, err := r.collection.UpdateOne(ctx, filter, updateDoc, opts)
	if err != nil {
		return fmt.Errorf("failed to create or update reaction record: %w", err)
	}

	// If a new document was inserted (upserted), update the ID and CreatedAt in the passed entity so the caller has the complete, persisted entity details.
	if res.UpsertedID != nil {
		if id, ok := res.UpsertedID.(uuid.UUID); ok {
			like.ID = id
		} else if idStr, ok := res.UpsertedID.(string); ok {
			if parsedID, parseErr := uuid.Parse(idStr); parseErr == nil {
				like.ID = parsedID
			}
		}
		like.CreatedAt = setOnInsertFields["created_at"].(time.Time)
	}

	return nil
}

// DeleteReaction marks a reaction record as deleted (soft delete) by its unique ID.
func (r *LikeRepository) DeleteReaction(ctx context.Context, reactionID uuid.UUID) error {
	filter := bson.M{"id": reactionID, "is_deleted": false}
	update := bson.M{"$set": bson.M{"is_deleted": true, "updated_at": time.Now()}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to delete reaction: %w", err)
	}
	if res.ModifiedCount == 0 {
		return errors.New("reaction not found")
	}
	return nil
}

// GetReactionByUserIDAndTargetID retrieves any active reaction (like or dislike) by a specific user on a specific target. Returns nil if no active reaction is found.
func (r *LikeRepository) GetReactionByUserIDAndTargetID(ctx context.Context, userID, targetID uuid.UUID) (*entity.Like, error) {
	var like entity.Like
	// Filter for active reactions
	filter := bson.M{"user_id": userID, "target_id": targetID, "is_deleted": false}

	err := r.collection.FindOne(ctx, filter).Decode(&like)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("reaction not found")
		}
		return nil, fmt.Errorf("failed to retrieve reaction: %w", err)
	}
	return &like, nil
}

// GetReactionByUserIDTargetIDAndType retrieves a specific type of active reaction (like or dislike) by a user on a target. Returns nil if no matching active reaction is found.
func (r *LikeRepository) GetReactionByUserIDTargetIDAndType(ctx context.Context, userID, targetID uuid.UUID, reactionType entity.LikeType) (*entity.Like, error) {
	var like entity.Like
	filter := bson.M{
		"user_id":    userID,
		"target_id":  targetID,
		"type":       reactionType,
		"is_deleted": false,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&like)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("reaction not found")
		}
		return nil, fmt.Errorf("failed to retrieve specific reaction: %w", err)
	}
	return &like, nil
}

// CountLikesByTargetID counts the number of active 'like' reactions for a specific target.
func (r *LikeRepository) CountLikesByTargetID(ctx context.Context, targetID uuid.UUID) (int64, error) {
	// Filter to count only active 'likes'
	filter := bson.M{"target_id": targetID, "type": entity.LIKE_TYPE_LIKE, "is_deleted": false}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count active likes: %w", err)
	}
	return count, nil
}

// CountDislikesByTargetID counts the number of active 'dislike' reactions for a specific target.
func (r *LikeRepository) CountDislikesByTargetID(ctx context.Context, targetID uuid.UUID) (int64, error) {
	// Filter to count only active 'dislikes'
	filter := bson.M{"target_id": targetID, "type": entity.LIKE_TYPE_DISLIKE, "is_deleted": false}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count active dislikes: %w", err)
	}
	return count, nil
}
