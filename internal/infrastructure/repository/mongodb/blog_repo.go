package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/contract"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
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

// buildBaseMatchFilter creates a bson.M filter from BlogFilterOptions.
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

// getSortOrder is a helper to convert sort order string to an integer.
func getSortOrder(sortOrder string) int {
	if sortOrder == "asc" {
		return 1
	}
	return -1
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

// GetBlogBySlug retrieves a single blog post by its unique slug.
func (r *BlogRepository) GetBlogBySlug(ctx context.Context, slug string) (*entity.Blog, error) {
	var blog entity.Blog
	filter := bson.M{"slug": slug, "is_deleted": false}

	err := r.collection.FindOne(ctx, filter).Decode(&blog)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("blog with slug '%s' not found or has been deleted: %w", slug, err)
		}
		return nil, fmt.Errorf("failed to retrieve blog post: %w", err)
	}

	return &blog, nil
}

// GetBlogs retrieves a list of blog posts with filtering, sorting, and pagination options.
func (r *BlogRepository) GetBlogs(ctx context.Context, opts *contract.BlogFilterOptions) ([]*entity.Blog, int64, error) {
	baseMatch := buildBaseMatchFilter(opts)

	// Base pipeline stages for filtering
	pipeline := []bson.M{
		{"$match": baseMatch},
	}

	// Conditionally add stages for tag filtering
	if len(opts.TagIDs) > 0 {
		pipeline = append(pipeline,
			bson.M{
				"$lookup": bson.M{
					"from":         "blog_tags",
					"localField":   "_id",
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

	// Define the sub-pipelines for $facet.
	sortField := "created_at"
	if opts.SortBy != "" {
		sortField = opts.SortBy
	}

	// The full pipeline includes all filters, followed by the $facet stage.
	fullPipeline := append(pipeline,
		bson.M{
			"$facet": bson.M{
				"totalCount": bson.A{
					bson.M{"$count": "total"},
				},
				"blogs": bson.A{
					bson.M{"$sort": bson.M{sortField: getSortOrder(opts.SortOrder)}},
					bson.M{"$skip": int64((opts.Page - 1) * opts.PageSize)},
					bson.M{"$limit": int64(opts.PageSize)},
				},
			},
		},
	)

	cursor, err := r.collection.Aggregate(ctx, fullPipeline)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to aggregate blogs: %w", err)
	}
	defer cursor.Close(ctx)

	var facetResults []struct {
		Blogs      []*entity.Blog `bson:"blogs"`
		TotalCount []struct {
			Total int64 `bson:"total"`
		} `bson:"totalCount"`
	}

	if err = cursor.All(ctx, &facetResults); err != nil {
		return nil, 0, fmt.Errorf("failed to decode facet results: %w", err)
	}

	if len(facetResults) == 0 {
		return []*entity.Blog{}, 0, nil
	}

	total := int64(0)
	if len(facetResults[0].TotalCount) > 0 {
		total = facetResults[0].TotalCount[0].Total
	}

	return facetResults[0].Blogs, total, nil
}

func (r *BlogRepository) UpdateBlog(ctx context.Context, blogID string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	filter := bson.M{"_id": blogID, "is_deleted": false}
	update := bson.M{"$set": updates}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update blog post: %w", err)
	}
	if res.ModifiedCount == 0 {
		var blog entity.Blog
		err := r.collection.FindOne(ctx, bson.M{"_id": blogID}).Decode(&blog)
		if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("blog post with ID %s not found", blogID)
		}
		return fmt.Errorf("blog post with ID %s was not modified (no new data to apply)", blogID)
	}

	return nil
}

func (r *BlogRepository) DeleteBlog(ctx context.Context, blogID string) error {
	filter := bson.M{"_id": blogID, "is_deleted": false}
	update := bson.M{"$set": bson.M{"is_deleted": true, "updated_at": time.Now()}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to delete blog post: %w", err)
	}
	if res.ModifiedCount == 0 {
		var blog entity.Blog
		err := r.collection.FindOne(ctx, bson.M{"_id": blogID}).Decode(&blog)
		if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("blog post with ID %s not found", blogID)
		}
		return fmt.Errorf("blog post with ID %s was not modified (possibly already deleted)", blogID)
	}

	return nil
}

// SearchBlogs searches for blog posts based on a query (title, author name, or author ID) and applies filter options.
func (r *BlogRepository) SearchBlogs(ctx context.Context, query string, opts *contract.BlogFilterOptions) ([]*entity.Blog, int64, error) {
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
	if _, err := uuid.Parse(query); err == nil {
		searchConditions = append(searchConditions, bson.M{"author_id": query})
	}

	// Aggregation pipeline for filtering and searching
	pipeline := []bson.M{
		{"$match": baseMatch},
		{"$lookup": bson.M{
			"from":         "users",
			"localField":   "author_id",
			"foreignField": "_id",
			"as":           "authorDetails",
		}},
		{"$unwind": "$authorDetails"},
		{"$match": bson.M{"$or": searchConditions}},
	}

	// If TagIDs are provided, add a $lookup to blog_tags and filter
	if len(opts.TagIDs) > 0 {
		pipeline = append(pipeline,
			bson.M{
				"$lookup": bson.M{
					"from":         "blog_tags",
					"localField":   "_id",
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

	// Determine the sort field, adding a prefix for joined collections if necessary.
	sortField := opts.SortBy
	switch sortField {
	case "":
		sortField = "created_at"
	case "username", "first_name", "last_name":
		sortField = "authorDetails." + sortField
	}

	// Use the $facet stage to perform counting and fetching results in a single pipeline.
	fullPipeline := append(pipeline,
		bson.M{
			"$facet": bson.M{
				"totalCount": bson.A{
					bson.M{"$count": "total"},
				},
				"blogs": bson.A{
					bson.M{"$sort": bson.M{sortField: getSortOrder(opts.SortOrder)}},
					bson.M{"$skip": int64((opts.Page - 1) * opts.PageSize)},
					bson.M{"$limit": int64(opts.PageSize)},
				},
			},
		},
	)

	cursor, err := r.collection.Aggregate(ctx, fullPipeline)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to aggregate search results: %w", err)
	}
	defer cursor.Close(ctx)

	var facetResults []struct {
		Blogs      []*entity.Blog `bson:"blogs"`
		TotalCount []struct {
			Total int64 `bson:"total"`
		} `bson:"totalCount"`
	}

	if err = cursor.All(ctx, &facetResults); err != nil {
		return nil, 0, fmt.Errorf("failed to decode search results from facet: %w", err)
	}

	if len(facetResults) == 0 {
		return []*entity.Blog{}, 0, nil
	}

	totalBlogs := int64(0)
	if len(facetResults[0].TotalCount) > 0 {
		totalBlogs = facetResults[0].TotalCount[0].Total
	}

	return facetResults[0].Blogs, totalBlogs, nil
}

// IncrementViewCount increments the view count of a specific blog post.
func (r *BlogRepository) IncrementViewCount(ctx context.Context, blogID string) error {
	filter := bson.M{"_id": blogID, "is_deleted": false}
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
func (r *BlogRepository) IncrementLikeCount(ctx context.Context, blogID string) error {
	filter := bson.M{"_id": blogID, "is_deleted": false}
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
func (r *BlogRepository) DecrementLikeCount(ctx context.Context, blogID string) error {
	filter := bson.M{"_id": blogID, "is_deleted": false}
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
func (r *BlogRepository) IncrementDislikeCount(ctx context.Context, blogID string) error {
	filter := bson.M{"_id": blogID, "is_deleted": false}
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
func (r *BlogRepository) DecrementDislikeCount(ctx context.Context, blogID string) error {
	filter := bson.M{"_id": blogID, "is_deleted": false}
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
func (r *BlogRepository) IncrementCommentCount(ctx context.Context, blogID string) error {
	filter := bson.M{"_id": blogID, "is_deleted": false}
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
func (r *BlogRepository) DecrementCommentCount(ctx context.Context, blogID string) error {
	filter := bson.M{"_id": blogID, "is_deleted": false}
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
func (r *BlogRepository) GetBlogCounts(ctx context.Context, blogID string) (viewCount, likeCount, dislikeCount, commentCount int, err error) {
	var blog entity.Blog
	filter := bson.M{"_id": blogID, "is_deleted": false}
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
func (r *BlogRepository) AddTagsToBlog(ctx context.Context, blogSlug string, tagIDs []string) error {
	if len(tagIDs) == 0 {
		return nil
	}

	// Check if the blog exists and is not deleted using its slug
	blog, err := r.GetBlogBySlug(ctx, blogSlug)
	if err != nil {
		return fmt.Errorf("blog not found: %w", err)
	}

	// Prepare documents for bulk insert
	var blogTags []interface{}
	for _, tagIDStr := range tagIDs {
		blogTag := entity.BlogTag{
			BlogID: blog.ID,
			TagID:  tagIDStr,
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
					fmt.Printf("Warning: Duplicate blog-tag association for blog %s and tag with index %d. Error: %v\n", blog.ID, e.Index, e)
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
func (r *BlogRepository) RemoveTagsFromBlog(ctx context.Context, blogSlug string, tagIDs []string) error {
	if len(tagIDs) == 0 {
		return nil
	}

	// Check if the blog exists and is not deleted using its slug
	blog, err := r.GetBlogBySlug(ctx, blogSlug)
	if err != nil {
		return fmt.Errorf("blog not found: %w", err)
	}

	// Prepare filter for deletion
	filter := bson.M{
		"blog_id": blog.ID,
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
