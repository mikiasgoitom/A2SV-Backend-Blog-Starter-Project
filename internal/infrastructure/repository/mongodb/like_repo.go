package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

// CreateLike inserts a new like record into the database.
func (r *LikeRepository) CreateLike(ctx context.Context, like *entity.Like) error {
	like.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, like)
	if err != nil {
		var writeException mongo.WriteException
		if errors.As(err, &writeException) {
			for _, e := range writeException.WriteErrors {
				if e.Code == 11000 {
					return errors.New("like already exists for this user and target")
				}
			}
		}
		return errors.New("failed to create like")
	}
	return nil
}

// DeleteLike deletes a like record from the database by its ID.
func (r *LikeRepository) DeleteLike(ctx context.Context, likeID string) error {
	filter := bson.M{"id": likeID}

	res, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return errors.New("failed to delete like")
	}
	if res.DeletedCount == 0 {
		return errors.New("like not found")
	}
	return nil
}

// GetLikeByUserIDAndTargetID retrieves a like record by the user ID and target ID.
func (r *LikeRepository) GetLikeByUserIDAndTargetID(ctx context.Context, userID, targetID string) (*entity.Like, error) {
	var like entity.Like
	filter := bson.M{"user_id": userID, "target_id": targetID}

	err := r.collection.FindOne(ctx, filter).Decode(&like)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("like not found")
		}
		return nil, errors.New("failed to retrieve like")
	}
	return &like, nil
}

// CountLikesByTargetID counts the number of likes for a specific target (blog or comment).
func (r *LikeRepository) CountLikesByTargetID(ctx context.Context, targetID string) (int64, error) {
	filter := bson.M{"target_id": targetID}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, errors.New("failed to count likes")
	}
	return count, nil
}
