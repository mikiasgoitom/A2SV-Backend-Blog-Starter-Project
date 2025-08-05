package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/contract"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BlogRepository represents the MongoDB implementation of the BlogRepository interface.
type BlogRepository struct {
	collection         *mongo.Collection // For blog posts
	blogTagsCollection *mongo.Collection // For blog-tag relationships
	usersCollection    *mongo.Collection // For accessing user data for search
}

// NewBlogRepository creates and returns a new BlogRepository instance.
func NewBlogRepository(db *mongo.Database) *BlogRepository {
	return &BlogRepository{
		collection:         db.Collection("blogs"),
		blogTagsCollection: db.Collection("blog_tags"),
		usersCollection:    db.Collection("users"),
	}
}

// CreateBlog inserts a new blog post record into the database.
func (r *BlogRepository) CreateBlog(ctx context.Context, blog *entity.Blog) error {
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, blog)
	if err != nil {
		return errors.New("failed to create blog post")
	}
	return nil
}

// GetBlogByID retrieves a single blog post by its unique ID.
func (r *BlogRepository) GetBlogByID(ctx context.Context, blogID uuid.UUID) (*entity.Blog, error) {
	var blog entity.Blog
	filter := bson.M{"id": blogID, "isDeleted": false}

	err := r.collection.FindOne(ctx, filter).Decode(&blog)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("blog not found")
		}
		return nil, errors.New("failed to retrieve blog post")
	}

	return &blog, nil
}

// GetBlogs retrieves a list of blog posts based on filtering, pagination, and sorting options.
func (r *BlogRepository) GetBlogs(ctx context.Context, opts *contract.BlogFilterOptions) ([]*entity.Blog, int64, error) {
	filter := bson.M{"isDeleted": false}
	findOptions := options.Find()

	if opts.Page > 0 && opts.PageSize > 0 {
		findOptions.SetSkip(int64((opts.Page - 1) * opts.PageSize))
		findOptions.SetLimit(int64(opts.PageSize))
	}

	if opts.SortBy != "" {
		sortOrder := -1
		if opts.SortOrder == "asc" {
			sortOrder = 1
		}
		findOptions.SetSort(bson.M{opts.SortBy: sortOrder})
	}

	if opts.DateFrom != nil {
		filter["createdAt"] = bson.M{"$gte": opts.DateFrom}
	}
	if opts.DateTo != nil {
		if filter["createdAt"] == nil {
			filter["createdAt"] = bson.M{}
		}
		filter["createdAt"].(bson.M)["$lte"] = opts.DateTo
	}

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, errors.New("failed to retrieve blog posts")
	}
	defer cursor.Close(ctx)

	var blogs []*entity.Blog
	if err = cursor.All(ctx, &blogs); err != nil {
		return nil, 0, errors.New("failed to decode blog posts")
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, errors.New("failed to count blog posts")
	}

	return blogs, total, nil
}

// UpdateBlog updates the details of an existing blog post by its ID.
func (r *BlogRepository) UpdateBlog(ctx context.Context, blogID uuid.UUID, updates map[string]interface{}) error {
	updates["updatedAt"] = time.Now()
	filter := bson.M{"id": blogID, "isDeleted": false}
	update := bson.M{"$set": updates}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New("failed to update blog post")
	}
	if res.ModifiedCount == 0 {
		return errors.New("blog post not found or no changes made")
	}

	return nil
}

// DeleteBlog marks a blog post as deleted by its ID.
func (r *BlogRepository) DeleteBlog(ctx context.Context, blogID uuid.UUID) error {
	filter := bson.M{"id": blogID}
	update := bson.M{"$set": bson.M{"isDeleted": true, "updatedAt": time.Now()}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New("failed to delete blog post")
	}
	if res.ModifiedCount == 0 {
		return errors.New("blog post not found")
	}

	return nil
}

// SearchBlogs searches for blog posts based on a query (title, author name, or author ID) and applies filter options.
func (r *BlogRepository) SearchBlogs(ctx context.Context, query string, opts *contract.BlogFilterOptions) ([]*entity.Blog, int64, error) {
	// Base match for non-deleted blogs
	baseMatch := bson.M{"isDeleted": false}

	// Add date filter to the base match if present
	if opts.DateFrom != nil {
		baseMatch["createdAt"] = bson.M{"$gte": opts.DateFrom}
	}
	if opts.DateTo != nil {
		if baseMatch["createdAt"] == nil {
			baseMatch["createdAt"] = bson.M{}
		}
		baseMatch["createdAt"].(bson.M)["$lte"] = opts.DateTo
	}

	// Define search conditions for title and author details
	searchConditions := []bson.M{
		{"title": bson.M{"$regex": query, "$options": "i"}},
	}

	// Add author name search conditions
	searchConditions = append(searchConditions,
		bson.M{"authorDetails.username": bson.M{"$regex": query, "$options": "i"}},
		bson.M{"authorDetails.firstName": bson.M{"$regex": query, "$options": "i"}},
		bson.M{"authorDetails.lastName": bson.M{"$regex": query, "$options": "i"}},
	)

	// If the query string is a valid UUID, also search by authorId
	if parsedUUID, err := uuid.Parse(query); err == nil {
		searchConditions = append(searchConditions, bson.M{"authorId": parsedUUID})
	}

	// Aggregation pipeline
	pipeline := []bson.M{
		// Initial match for non-deleted blogs and date filter
		{"$match": baseMatch},
		// Look up author details from the 'users' collection
		{"$lookup": bson.M{
			"from":         "users",
			"localField":   "authorId",
			"foreignField": "id",
			"as":           "authorDetails",
		}},
		// Unwind the authorDetails array (each blog will now have an authorDetails object)
		{"$unwind": "$authorDetails"},
		// Match based on the search query (title OR author name/ID)
		{"$match": bson.M{"$or": searchConditions}},
	}

	// Create a separate pipeline for counting, by appending $count stage
	countPipeline := append([]bson.M{}, pipeline...)
	countPipeline = append(countPipeline, bson.M{"$count": "total"})

	countCursor, err := r.collection.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}
	defer countCursor.Close(ctx)

	var countResult []struct {
		Total int64 `bson:"total"`
	}
	if err = countCursor.All(ctx, &countResult); err != nil {
		return nil, 0, fmt.Errorf("failed to decode count result: %w", err)
	}

	totalBlogs := int64(0)
	if len(countResult) > 0 {
		totalBlogs = countResult[0].Total
	}

	// --- Apply sorting to the main pipeline ---
	if opts.SortBy != "" {
		sortOrder := 1
		if opts.SortOrder == "desc" {
			sortOrder = -1
		}
		pipeline = append(pipeline, bson.M{"$sort": bson.M{opts.SortBy: sortOrder}})
	} else {
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"createdAt": -1}})
	}

	// --- Apply pagination to the main pipeline ---
	if opts.Page > 0 && opts.PageSize > 0 {
		pipeline = append(pipeline, bson.M{"$skip": int64((opts.Page - 1) * opts.PageSize)})
		pipeline = append(pipeline, bson.M{"$limit": int64(opts.PageSize)})
	}

	pipeline = append(pipeline, bson.M{"$replaceRoot": bson.M{"newRoot": "$blogDetails"}})

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search for blog posts: %w", err)
	}
	defer cursor.Close(ctx)

	var blogs []*entity.Blog
	if err = cursor.All(ctx, &blogs); err != nil {
		return nil, 0, fmt.Errorf("failed to decode search results: %w", err)
	}

	return blogs, totalBlogs, nil
}

// IncrementViewCount increments the view count of a specific blog post.
func (r *BlogRepository) IncrementViewCount(ctx context.Context, blogID uuid.UUID) error {
	filter := bson.M{"id": blogID}
	update := bson.M{"$inc": bson.M{"viewCount": 1}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New("failed to increment view count")
	}
	if res.ModifiedCount == 0 {
		return errors.New("blog post not found")
	}

	return nil
}

// IncrementLikeCount increments the like count of a specific blog post.
func (r *BlogRepository) IncrementLikeCount(ctx context.Context, blogID uuid.UUID) error {
	filter := bson.M{"id": blogID}
	update := bson.M{"$inc": bson.M{"likeCount": 1}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New("failed to increment like count")
	}
	if res.ModifiedCount == 0 {
		return errors.New("blog post not found")
	}

	return nil
}

// DecrementLikeCount decrements the like count of a specific blog post.
func (r *BlogRepository) DecrementLikeCount(ctx context.Context, blogID uuid.UUID) error {
	filter := bson.M{"id": blogID}
	update := bson.M{"$inc": bson.M{"likeCount": -1}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New("failed to decrement like count")
	}
	if res.ModifiedCount == 0 {
		return errors.New("blog post not found")
	}

	return nil
}

// IncrementDislikeCount increments the dislike count of a specific blog post.
func (r *BlogRepository) IncrementDislikeCount(ctx context.Context, blogID uuid.UUID) error {
	filter := bson.M{"id": blogID}
	update := bson.M{"$inc": bson.M{"dislikeCount": 1}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New("failed to increment dislike count")
	}
	if res.ModifiedCount == 0 {
		return errors.New("blog post not found")
	}

	return nil
}

// DecrementDislikeCount decrements the dislike count of a specific blog post.
func (r *BlogRepository) DecrementDislikeCount(ctx context.Context, blogID uuid.UUID) error {
	filter := bson.M{"id": blogID}
	update := bson.M{"$inc": bson.M{"dislikeCount": -1}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New("failed to decrement dislike count")
	}
	if res.ModifiedCount == 0 {
		return errors.New("blog post not found")
	}

	return nil
}

// IncrementCommentCount increments the comment count of a specific blog post.
func (r *BlogRepository) IncrementCommentCount(ctx context.Context, blogID uuid.UUID) error {
	filter := bson.M{"id": blogID}
	update := bson.M{"$inc": bson.M{"commentCount": 1}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New("failed to increment comment count")
	}
	if res.ModifiedCount == 0 {
		return errors.New("blog post not found")
	}

	return nil
}

// DecrementCommentCount decrements the comment count of a specific blog post.
func (r *BlogRepository) DecrementCommentCount(ctx context.Context, blogID uuid.UUID) error {
	filter := bson.M{"id": blogID}
	update := bson.M{"$inc": bson.M{"commentCount": -1}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New("failed to decrement comment count")
	}
	if res.ModifiedCount == 0 {
		return errors.New("blog post not found")
	}

	return nil
}

// AddUserLike adds a like or dislike record for a user on a blog post.
func (r *BlogRepository) AddUserLike(ctx context.Context, blogID, userID uuid.UUID, likeType string) error {
	// Check if user already liked/disliked
	filter := bson.M{"blogId": blogID, "userId": userID}
	update := bson.M{
		"$set": bson.M{
			"type":   likeType,
			"status": "active",
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.Database().Collection("blog_likes").UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return errors.New("failed to add user like")
	}
	return nil
}

// RemoveUserLike removes a like or dislike record for a user on a blog post.
func (r *BlogRepository) RemoveUserLike(ctx context.Context, blogID, userID uuid.UUID) error {
	filter := bson.M{"blogId": blogID, "userId": userID}
	_, err := r.collection.Database().Collection("blog_likes").DeleteOne(ctx, filter)
	if err != nil {
		return errors.New("failed to remove user like")
	}
	return nil
}

// HasUserLiked checks if a user has liked or disliked a blog post.
func (r *BlogRepository) HasUserLiked(ctx context.Context, blogID, userID uuid.UUID) (string, bool, error) {
	filter := bson.M{"blogId": blogID, "userId": userID}
	var result struct {
		Type   string `bson:"type"`
		Status string `bson:"status"`
	}
	err := r.collection.Database().Collection("blog_likes").FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", false, nil
		}
		return "", false, errors.New("failed to check user like")
	}
	return result.Type, result.Status == "active", nil
}

// GetBlogCounts returns the current counts for a blog post.
func (r *BlogRepository) GetBlogCounts(ctx context.Context, blogID uuid.UUID) (viewCount, likeCount, dislikeCount, commentCount int, err error) {
	var blog entity.Blog
	filter := bson.M{"id": blogID}
	err = r.collection.FindOne(ctx, filter).Decode(&blog)
	if err != nil {
		return 0, 0, 0, 0, errors.New("failed to get blog counts")
	}
	return blog.ViewCount, blog.LikeCount, blog.DislikeCount, blog.CommentCount, nil
}

// AddTagsToBlog associates one or more tags with a blog post.
func (r *BlogRepository) AddTagsToBlog(ctx context.Context, blogID uuid.UUID, tagIDs []uuid.UUID) error {
	if len(tagIDs) == 0 {
		return nil
	}

	// Check if the blog exists and is not deleted
	_, err := r.GetBlogByID(ctx, blogID)
	if err != nil {
		return fmt.Errorf("blog not found: %w", err)
	}

	// Prepare documents for bulk insert
	var blogTags []interface{}
	for _, tagID := range tagIDs {
		blogTag := entity.BlogTag{
			BlogID: blogID,
			TagID:  tagID,
		}
		blogTags = append(blogTags, blogTag)
	}

	// Insert many, ignoring duplicate key errors if a blog-tag pair already exists
	opts := options.InsertMany().SetOrdered(false)
	_, err = r.blogTagsCollection.InsertMany(ctx, blogTags, opts)
	if err != nil {
		// Check for duplicate key errors specifically
		if writeException, ok := err.(mongo.BulkWriteException); ok {
			for _, e := range writeException.WriteErrors {
				if e.Code == 11000 {
					// Log or ignore duplicate errors, as they mean the tag was already associated
					fmt.Printf("Warning: Duplicate blog-tag association for blog %s and tag %s. Error: %v\n", blogID, e.Raw, e)
				} else {
					return fmt.Errorf("failed to add tags: %w", err)
				}
			}
		} else {
			return fmt.Errorf("failed to add tags: %w", err)
		}
	}
	return nil
}

// RemoveTagsFromBlog disassociates one or more tags from a blog post.
func (r *BlogRepository) RemoveTagsFromBlog(ctx context.Context, blogID uuid.UUID, tagIDs []uuid.UUID) error {
	if len(tagIDs) == 0 {
		return nil
	}

	// Check if the blog exists and is not deleted
	_, err := r.GetBlogByID(ctx, blogID)
	if err != nil {
		return fmt.Errorf("blog not found: %w", err)
	}

	// Prepare filter for deletion
	filter := bson.M{
		"blogId": blogID,
		"tagId":  bson.M{"$in": tagIDs},
	}

	res, err := r.blogTagsCollection.DeleteMany(ctx, filter)
	if err != nil {
		return errors.New("failed to remove tags from blog")
	}
	if res.DeletedCount == 0 {
		return errors.New("no matching tags found to remove for the blog post")
	}

	return nil
}

// GetBlogsByTagID retrieves a list of blog posts associated with a specific tag ID, applying pagination and sorting options.
func (r *BlogRepository) GetBlogsByTagID(ctx context.Context, tagID uuid.UUID, opts *contract.BlogFilterOptions) ([]*entity.Blog, int64, error) {
	// Aggregation pipeline to join blog_tags with blogs
	pipeline := []bson.M{
		// Match blog_tags documents by the given tagId.
		{"$match": bson.M{"tagId": tagID}},
		// Look up the corresponding blog documents from the 'blogs' collection
		{"$lookup": bson.M{
			"from":         "blogs",
			"localField":   "blogId",
			"foreignField": "id",
			"as":           "blogDetails",
		}},
		// Unwind the blogDetails array (each blog_tag document will now have a blogDetails object)
		{"$unwind": "$blogDetails"},
		// Match only active (not deleted) blogs
		{"$match": bson.M{"blogDetails.isDeleted": false}},
	}

	// Add sorting
	if opts.SortBy != "" {
		sortOrder := 1
		if opts.SortOrder == "desc" {
			sortOrder = -1
		}
		pipeline = append(pipeline, bson.M{"$sort": bson.M{fmt.Sprintf("blogDetails.%s", opts.SortBy): sortOrder}})
	} else {
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"blogDetails.createdAt": -1}})
	}

	// Count total documents before pagination
	countPipeline := append(pipeline, bson.M{"$count": "total"})
	countCursor, err := r.blogTagsCollection.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, 0, errors.New("failed to count blogs by tag")
	}
	defer countCursor.Close(ctx)

	var countResult []struct {
		Total int64 `bson:"total"`
	}
	if err = countCursor.All(ctx, &countResult); err != nil {
		return nil, 0, errors.New("failed to decode count result")
	}

	totalBlogs := int64(0)
	if len(countResult) > 0 {
		totalBlogs = countResult[0].Total
	}

	// Add pagination to the main pipeline
	if opts.Page > 0 && opts.PageSize > 0 {
		pipeline = append(pipeline, bson.M{"$skip": int64((opts.Page - 1) * opts.PageSize)})
		pipeline = append(pipeline, bson.M{"$limit": int64(opts.PageSize)})
	}

	// Project to shape the output as an entity.Blog
	pipeline = append(pipeline, bson.M{"$replaceRoot": bson.M{"newRoot": "$blogDetails"}})

	cursor, err := r.blogTagsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, errors.New("failed to retrieve blogs by tag")
	}
	defer cursor.Close(ctx)

	var blogs []*entity.Blog
	if err = cursor.All(ctx, &blogs); err != nil {
		return nil, 0, fmt.Errorf("failed to decode blogs by tag: %w", err)
	}

	return blogs, totalBlogs, nil
}
