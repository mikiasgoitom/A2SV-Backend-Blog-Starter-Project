package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
)

// CommentRepository represents the MongoDB implementation of the ICommentRepository interface.
type CommentRepository struct {
	collection *mongo.Collection
}

// NewCommentRepository creates and returns a new CommentRepository instance.
func NewCommentRepository(db *mongo.Database) *CommentRepository {
	return &CommentRepository{
		collection: db.Collection("comments"),
	}
}

// CreateComment inserts a new comment record into the database.
func (r *CommentRepository) CreateComment(ctx context.Context, comment *entity.Comment) error {
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, comment)
	if err != nil {
		return errors.New("failed to create comment")
	}
	return nil
}

// GetCommentByID retrieves a single comment by its unique ID.
func (r *CommentRepository) GetCommentByID(ctx context.Context, commentID uuid.UUID) (*entity.Comment, error) {
	var comment entity.Comment
	filter := bson.M{"id": commentID, "is_deleted": false}

	err := r.collection.FindOne(ctx, filter).Decode(&comment)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("comment not found")
		}
		return nil, errors.New("failed to retrieve comment")
	}
	return &comment, nil
}

// GetCommentsByBlogID retrieves a list of comments for a specific blog post, with pagination.
func (r *CommentRepository) GetCommentsByBlogID(ctx context.Context, blogID uuid.UUID, page, pageSize int) ([]*entity.Comment, int64, error) {
	filter := bson.M{"blog_id": blogID, "is_deleted": false}
	findOptions := options.Find()

	if page > 0 && pageSize > 0 {
		findOptions.SetSkip(int64((page - 1) * pageSize))
		findOptions.SetLimit(int64(pageSize))
	}

	findOptions.SetSort(bson.M{"created_at": 1})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, errors.New("failed to retrieve comments")
	}
	defer cursor.Close(ctx)

	var comments []*entity.Comment
	if err = cursor.All(ctx, &comments); err != nil {
		return nil, 0, errors.New("failed to decode comments")
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, errors.New("failed to count comments")
	}

	return comments, total, nil
}

// UpdateComment updates the details of an existing comment by its ID.
func (r *CommentRepository) UpdateComment(ctx context.Context, commentID uuid.UUID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	filter := bson.M{"id": commentID, "is_deleted": false}
	update := bson.M{"$set": updates}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New("failed to update comment")
	}
	if res.ModifiedCount == 0 {
		return errors.New("comment not found or no changes made")
	}
	return nil
}

// DeleteComment marks a comment as deleted by its ID (soft delete).
func (r *CommentRepository) DeleteComment(ctx context.Context, commentID uuid.UUID) error {
	filter := bson.M{"id": commentID}
	update := bson.M{"$set": bson.M{"is_deleted": true, "updated_at": time.Now()}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New("failed to delete comment")
	}
	if res.ModifiedCount == 0 {
		return errors.New("comment not found")
	}
	return nil
}
