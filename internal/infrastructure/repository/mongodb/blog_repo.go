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

// --- Helper Functions for code refactoring ---

// buildBaseMatchFilter creates a bson.M filter from BlogFilterOptions.
// All pointer options are now consistently dereferenced.
func buildBaseMatchFilter(opts *contract.BlogFilterOptions) bson.M {
	baseMatch := bson.M{"is_deleted": false}

	// Apply date range filter
	if opts.DateFrom != nil || opts.DateTo != nil {
		createdAtFilter := bson.M{}
		if opts.DateFrom != nil {
			createdAtFilter["$gte"] = *opts.DateFrom
		}
		if opts.DateTo != nil {
			createdAtFilter["$lte"] = *opts.DateTo
		}
		baseMatch["created_at"] = createdAtFilter
	}

	// Apply view count range filter
	if opts.MinViews != nil || opts.MaxViews != nil {
		viewCountFilter := bson.M{}
		if opts.MinViews != nil {
			viewCountFilter["$gte"] = *opts.MinViews
		}
		if opts.MaxViews != nil {
			viewCountFilter["$lte"] = *opts.MaxViews
		}
		baseMatch["view_count"] = viewCountFilter
	}

	// Apply like count range filter
	if opts.MinLikes != nil || opts.MaxLikes != nil {
		likeCountFilter := bson.M{}
		if opts.MinLikes != nil {
			likeCountFilter["$gte"] = *opts.MinLikes
		}
		if opts.MaxLikes != nil {
			likeCountFilter["$lte"] = *opts.MaxLikes
		}
		baseMatch["like_count"] = likeCountFilter
	}

	// Apply author ID filter
	if opts.AuthorID != nil {
		baseMatch["author_id"] = *opts.AuthorID
	}

	return baseMatch
}

// buildAggregationSortAndPaginationStages creates the sort and pagination stages for an aggregation pipeline.
// `fieldPrefix` is used for cases where sorting fields are nested (e.g., "blogDetails.").
func buildAggregationSortAndPaginationStages(opts *contract.BlogFilterOptions, fieldPrefix string) []bson.M {
	stages := []bson.M{}

	sortOrder := -1
	if opts.SortOrder == "asc" {
		sortOrder = 1
	}

	sortBy := "created_at"
	if opts.SortBy != "" {
		sortBy = opts.SortBy
	}

	sortField := fmt.Sprintf("%s%s", fieldPrefix, sortBy)
	stages = append(stages, bson.M{"$sort": bson.M{sortField: sortOrder}})

	if opts.Page > 0 && opts.PageSize > 0 {
		stages = append(stages, bson.M{"$skip": int64((opts.Page - 1) * opts.PageSize)})
		stages = append(stages, bson.M{"$limit": int64(opts.PageSize)})
	}

	return stages
}

// CreateBlog inserts a new blog post record into the database.
func (r *BlogRepository) CreateBlog(ctx context.Context, blog *entity.Blog) error {
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, blog)
	if err != nil {
		return fmt.Errorf("failed to create blog post: %w", err)
	}
	return nil
}

// GetBlogByID retrieves a single blog post by its unique ID.
func (r *BlogRepository) GetBlogByID(ctx context.Context, blogID uuid.UUID) (*entity.Blog, error) {
	var blog entity.Blog
	filter := bson.M{"id": blogID, "is_deleted": false}

	err := r.collection.FindOne(ctx, filter).Decode(&blog)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("blog not found")
		}
		return nil, fmt.Errorf("failed to retrieve blog post: %w", err)
	}

	return &blog, nil
}

// GetBlogs retrieves a list of blog posts with filtering, sorting, and pagination options.
// This method now uses a single aggregation pipeline for all filter scenarios.
func (r *BlogRepository) GetBlogs(ctx context.Context, opts *contract.BlogFilterOptions) ([]*entity.Blog, int64, error) {
	// Build the base filter using the helper function.
	baseMatch := buildBaseMatchFilter(opts)

	// A pipeline is always used to handle all filter combinations.
	pipeline := []bson.M{
		{"$match": baseMatch},
	}

	// Conditionally add stages for tag filtering
	if len(opts.TagIDs) > 0 {
		pipeline = append(pipeline,
			bson.M{
				"$lookup": bson.M{
					"from":         "blog_tags",
					"localField":   "id",
					"foreignField": "blog_id",
					"as":           "blogTags",
				},
			},
			bson.M{"$unwind": "$blogTags"},
			bson.M{"$match": bson.M{"blogTags.tag_id": bson.M{"$in": opts.TagIDs}}},
			// Use a group stage to get unique blogs after the unwind
			bson.M{
				"$group": bson.M{
					"_id":  "$id",
					"blog": bson.M{"$first": "$$ROOT"},
				},
			},
			// Replace the root with the original blog document
			bson.M{"$replaceRoot": bson.M{"newRoot": "$blog"}},
		)
	}

	// Create a separate pipeline for counting by appending the $count stage.
	countPipeline := append([]bson.M{}, pipeline...)
	countPipeline = append(countPipeline, bson.M{"$count": "total"})

	// Execute the count pipeline.
	countCursor, err := r.collection.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count blogs: %w", err)
	}
	defer countCursor.Close(ctx)

	var countResult []struct {
		Total int64 `bson:"total"`
	}
	if err = countCursor.All(ctx, &countResult); err != nil {
		return nil, 0, fmt.Errorf("failed to decode count result for blogs: %w", err)
	}

	total := int64(0)
	if len(countResult) > 0 {
		total = countResult[0].Total
	}

	// Add sorting and pagination to the main pipeline.
	pipeline = append(pipeline, buildAggregationSortAndPaginationStages(opts, "")...)

	// Execute the main pipeline.
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve blogs: %w", err)
	}
	defer cursor.Close(ctx)

	var blogs []*entity.Blog
	if err = cursor.All(ctx, &blogs); err != nil {
		return nil, 0, fmt.Errorf("failed to decode blogs: %w", err)
	}

	return blogs, total, nil
}

// UpdateBlog updates the details of an existing blog post by its ID.
func (r *BlogRepository) UpdateBlog(ctx context.Context, blogID uuid.UUID, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	filter := bson.M{"id": blogID, "is_deleted": false}
	update := bson.M{"$set": updates}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update blog post: %w", err)
	}
	if res.ModifiedCount == 0 {
		return errors.New("blog post not found or no changes made")
	}

	return nil
}

// DeleteBlog marks a blog post as deleted by its ID.
func (r *BlogRepository) DeleteBlog(ctx context.Context, blogID uuid.UUID) error {
	filter := bson.M{"id": blogID}
	update := bson.M{"$set": bson.M{"is_deleted": true, "updated_at": time.Now()}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to delete blog post: %w", err)
	}
	if res.ModifiedCount == 0 {
		return errors.New("blog post not found")
	}

	return nil
}

// SearchBlogs searches for blog posts based on a query (title, author name, or author ID) and applies filter options.
func (r *BlogRepository) SearchBlogs(ctx context.Context, query string, opts *contract.BlogFilterOptions) ([]*entity.Blog, int64, error) {
	// Build the base filter using the helper function.
	baseMatch := buildBaseMatchFilter(opts)

	// Define search conditions for title and author details
	searchConditions := []bson.M{
		{"title": bson.M{"$regex": query, "$options": "i"}},
	}

	// Add author name search conditions
	searchConditions = append(searchConditions,
		bson.M{"authorDetails.username": bson.M{"$regex": query, "$options": "i"}},
		bson.M{"authorDetails.first_name": bson.M{"$regex": query, "$options": "i"}},
		bson.M{"authorDetails.last_name": bson.M{"$regex": query, "$options": "i"}},
	)

	// If the query string is a valid UUID, also search by author_id
	if parsedUUID, err := uuid.Parse(query); err == nil {
		searchConditions = append(searchConditions, bson.M{"author_id": parsedUUID})
	}

	// Aggregation pipeline
	pipeline := []bson.M{
		// Initial match for non-deleted blogs and date filter
		{"$match": baseMatch},
		// Look up author details from the 'users' collection
		{"$lookup": bson.M{
			"from":         "users",
			"localField":   "author_id",
			"foreignField": "id",
			"as":           "authorDetails",
		}},
		// Unwind the authorDetails array (each blog will now have an authorDetails object)
		{"$unwind": "$authorDetails"},
		// Match based on the search query (title OR author name/ID)
		{"$match": bson.M{"$or": searchConditions}},
	}

	// If TagIDs are provided, add a $lookup to blog_tags and filter
	if len(opts.TagIDs) > 0 {
		pipeline = append(pipeline,
			bson.M{
				"$lookup": bson.M{
					"from":         "blog_tags",
					"localField":   "id",
					"foreignField": "blog_id",
					"as":           "blogTags",
				},
			},
			bson.M{"$unwind": "$blogTags"},
			bson.M{"$match": bson.M{"blogTags.tag_id": bson.M{"$in": opts.TagIDs}}},
			// Use a group stage to get unique blogs after the unwind
			bson.M{
				"$group": bson.M{
					"_id":  "$id",
					"blog": bson.M{"$first": "$$ROOT"},
				},
			},
			// Replace the root with the original blog document
			bson.M{"$replaceRoot": bson.M{"newRoot": "$blog"}},
		)
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

	// Add sorting and pagination to the main pipeline using the helper
	pipeline = append(pipeline, buildAggregationSortAndPaginationStages(opts, "")...)

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
	update := bson.M{"$inc": bson.M{"view_count": 1}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to increment view count: %w", err)
	}
	if res.ModifiedCount == 0 {
		return errors.New("blog post not found")
	}

	return nil
}

// IncrementLikeCount increments the like count of a specific blog post.
func (r *BlogRepository) IncrementLikeCount(ctx context.Context, blogID uuid.UUID) error {
	filter := bson.M{"id": blogID}
	update := bson.M{"$inc": bson.M{"like_count": 1}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to increment like count: %w", err)
	}
	if res.ModifiedCount == 0 {
		return errors.New("blog post not found")
	}

	return nil
}

// DecrementLikeCount decrements the like count of a specific blog post.
func (r *BlogRepository) DecrementLikeCount(ctx context.Context, blogID uuid.UUID) error {
	filter := bson.M{"id": blogID}
	update := bson.M{"$inc": bson.M{"like_count": -1}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to decrement like count: %w", err)
	}
	if res.ModifiedCount == 0 {
		return errors.New("blog post not found")
	}

	return nil
}

// IncrementDislikeCount increments the dislike count of a specific blog post.
func (r *BlogRepository) IncrementDislikeCount(ctx context.Context, blogID uuid.UUID) error {
	filter := bson.M{"id": blogID}
	update := bson.M{"$inc": bson.M{"dislike_count": 1}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to increment dislike count: %w", err)
	}
	if res.ModifiedCount == 0 {
		return errors.New("blog post not found")
	}

	return nil
}

// DecrementDislikeCount decrements the dislike count of a specific blog post.
func (r *BlogRepository) DecrementDislikeCount(ctx context.Context, blogID uuid.UUID) error {
	filter := bson.M{"id": blogID}
	update := bson.M{"$inc": bson.M{"dislike_count": -1}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to decrement dislike count: %w", err)
	}
	if res.ModifiedCount == 0 {
		return errors.New("blog post not found")
	}

	return nil
}

// IncrementCommentCount increments the comment count of a specific blog post.
func (r *BlogRepository) IncrementCommentCount(ctx context.Context, blogID uuid.UUID) error {
	filter := bson.M{"id": blogID}
	update := bson.M{"$inc": bson.M{"comment_count": 1}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to increment comment count: %w", err)
	}
	if res.ModifiedCount == 0 {
		return errors.New("blog post not found")
	}

	return nil
}

// DecrementCommentCount decrements the comment count of a specific blog post.
func (r *BlogRepository) DecrementCommentCount(ctx context.Context, blogID uuid.UUID) error {
	filter := bson.M{"id": blogID}
	update := bson.M{"$inc": bson.M{"comment_count": -1}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to decrement comment count: %w", err)
	}
	if res.ModifiedCount == 0 {
		return errors.New("blog post not found")
	}

	return nil
}

// GetBlogCounts returns the current counts for a blog post.
func (r *BlogRepository) GetBlogCounts(ctx context.Context, blogID uuid.UUID) (viewCount, likeCount, dislikeCount, commentCount int, err error) {
	var blog entity.Blog
	filter := bson.M{"id": blogID}
	err = r.collection.FindOne(ctx, filter).Decode(&blog)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, 0, 0, 0, errors.New("blog not found")
		}
		return 0, 0, 0, 0, fmt.Errorf("failed to get blog counts: %w", err)
	}
	return blog.ViewCount, blog.LikeCount, blog.DislikeCount, blog.CommentCount, nil
}

// AddTagsToBlog associates one or more tags with a blog post.
func (r *BlogRepository) AddTagsToBlog(ctx context.Context, blogID uuid.UUID, tagIDs []uuid.UUID) error {
	if len(tagIDs) == 0 {
		return nil
	}

	_, err := r.GetBlogByID(ctx, blogID)
	if err != nil {
		return fmt.Errorf("blog not found: %w", err)
	}

	var blogTags []interface{}
	for _, tagID := range tagIDs {
		blogTag := entity.BlogTag{
			BlogID: blogID,
			TagID:  tagID,
		}
		blogTags = append(blogTags, blogTag)
	}

	opts := options.InsertMany().SetOrdered(false)
	_, err = r.blogTagsCollection.InsertMany(ctx, blogTags, opts)
	if err != nil {
		if writeException, ok := err.(mongo.BulkWriteException); ok {
			type tempBlogTag struct {
				TagID uuid.UUID `bson:"tag_id"`
			}
			for _, e := range writeException.WriteErrors {
				if e.Code == 11000 {
					var failedTagID uuid.UUID
					var tempTag tempBlogTag
					if unmarshalErr := bson.Unmarshal(e.Raw, &tempTag); unmarshalErr == nil {
						failedTagID = tempTag.TagID
					} else {
						failedTagID = uuid.Nil
					}
					fmt.Printf("Warning: Duplicate blog-tag association for blog %s and tag %s. Error: %v\n", blogID, failedTagID, e)
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

	_, err := r.GetBlogByID(ctx, blogID)
	if err != nil {
		return fmt.Errorf("blog not found: %w", err)
	}

	filter := bson.M{
		"blog_id": blogID,
		"tag_id":  bson.M{"$in": tagIDs},
	}

	res, err := r.blogTagsCollection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to remove tags from blog: %w", err)
	}
	if res.DeletedCount == 0 {
		return errors.New("no matching tags found to remove for the blog post")
	}

	return nil
}
