package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MediaRepository represents the MongoDB implementation of the IMediaRepository interface.
type MediaRepository struct {
	collection *mongo.Collection
}

// NewMediaRepository creates and returns a new MediaRepository instance.
func NewMediaRepository(db *mongo.Database) *MediaRepository {
	return &MediaRepository{
		collection: db.Collection("media"),
	}
}

// CreateMedia inserts a new media record into the database.
func (r *MediaRepository) CreateMedia(ctx context.Context, media *entity.Media) error {
	media.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, media)
	if err != nil {
		return errors.New("failed to create media record")
	}
	return nil
}

// GetMediaByID retrieves a single media record by its unique ID.
func (r *MediaRepository) GetMediaByID(ctx context.Context, mediaID uuid.UUID) (*entity.Media, error) {
	var media entity.Media
	filter := bson.M{"id": mediaID}

	err := r.collection.FindOne(ctx, filter).Decode(&media)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("media not found")
		}
		return nil, errors.New("failed to retrieve media record")
	}
	return &media, nil
}

// DeleteMedia deletes a media record by its ID.
func (r *MediaRepository) DeleteMedia(ctx context.Context, mediaID uuid.UUID) error {
	filter := bson.M{"id": mediaID}
	res, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return errors.New("failed to delete media record")
	}
	if res.DeletedCount == 0 {
		return errors.New("media record not found")
	}
	return nil
}
