package mongodb

import (
	"context"
	"fmt"

	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/contract"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
)

// userAvatarRepository implements IUserAvatarRepository for MongoDB
type userAvatarRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewUserAvatarRepository creates a new MongoDB implementation of IUserAvatarRepository
func NewUserAvatarRepository(db *mongo.Database) contract.IUserAvatarRepository {
	return &userAvatarRepository{
		db:         db,
		collection: db.Collection("users"),
	}
}

// CreateAvatarURI saves a new avatar URI for a user
func (r *userAvatarRepository) CreateAvatarURI(ctx context.Context, userID, uri string) error {
	fmt.Printf("Debug: userID=%s\n", userID)

	filter := bson.M{"_id": userID}
	update := bson.M{"$set": bson.M{"avatar_uri": uri}}

	fmt.Printf("Debug: filter=%v, update=%v\n", filter, update)

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Printf("Debug: UpdateOne error=%v\n", err)
		return fmt.Errorf("failed to create avatar URI: %w", err)
	}

	fmt.Printf("Debug: MatchedCount=%d\n", result.MatchedCount)

	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// ReadAvatarURI retrieves the avatar URI from a user's record
func (r *userAvatarRepository) ReadAvatarURI(ctx context.Context, userID string) (string, error) {
	var user struct {
		AvatarURI string `bson:"avatar_uri"`
	}

	filter := bson.M{"_id": userID}
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", nil
		}
		return "", fmt.Errorf("failed to read avatar URI: %w", err)
	}

	return user.AvatarURI, nil
}

// UpdateAvatarURI updates the user's avatar URI with a new one
func (r *userAvatarRepository) UpdateAvatarURI(ctx context.Context, userID, newURI string) error {
	filter := bson.M{"_id": userID}
	update := bson.M{"$set": bson.M{"avatar_uri": newURI}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update avatar URI: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// DeleteAvatarURI removes the avatar URI from a user's record
func (r *userAvatarRepository) DeleteAvatarURI(ctx context.Context, userID string) error {
	filter := bson.M{"_id": userID}
	update := bson.M{"$unset": bson.M{"avatar_uri": ""}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to delete avatar URI: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// ReadMediaByURI retrieves media metadata using the avatar's URI
func (r *userAvatarRepository) ReadMediaByURI(ctx context.Context, uri string) (*entity.Media, error) {
	// This would typically query a media collection
	// For now, we'll return nil as the existing IMediaRepository handles this
	return nil, nil
}
