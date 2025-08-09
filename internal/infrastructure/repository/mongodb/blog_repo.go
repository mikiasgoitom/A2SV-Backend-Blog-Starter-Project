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
	collection          *mongo.Collection // For blog posts
	usersCollection     *mongo.Collection // For accessing user data for search
	blogViewsCollection *mongo.Collection // For tracking blog views
}

// NewBlogRepository creates and returns a new BlogRepository instance.
func NewBlogRepository(db *mongo.Database) *BlogRepository {
	return &BlogRepository{
		collection:          db.Collection("blogs"),
		usersCollection:     db.Collection("users"),
		blogViewsCollection: db.Collection("blog_views"),
	}
}

// getSortOrder is a helper to convert sort order string to an integer.
func getSortOrder(sortOrder string) int {
	if sortOrder == "asc" {
		return 1
	}
	return -1
}

// buildBaseMatchFilter creates a bson.M filter from BlogFilterOptions.
func buildBaseMatchFilter(opts *contract.BlogFilterOptions) bson.M {
	baseMatch := bson.M{"is_deleted": false}

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

	if opts.AuthorID != nil {
		baseMatch["author_id"] = *opts.AuthorID
	}

	// Updated to handle embedded tags from BlogFilterOptions
	if len(opts.TagIDs) > 0 {
		baseMatch["tags"] = bson.M{"$in": opts.TagIDs}
	}

	return baseMatch
}

// CreateBlog inserts a new blog post record into the database.
func (r *BlogRepository) CreateBlog(ctx context.Context, blog *entity.Blog) error {
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()
	if blog.Tags == nil {
		blog.Tags = []string{} // Ensure tags is not nil to avoid DB errors
	}
	_, err := r.collection.InsertOne(ctx, blog)
	if err != nil {
		return fmt.Errorf("failed to create blog post: %w", err)
	}
	return nil
}

// GetBlogByID retrieves a single blog post by its unique id.
func (r *BlogRepository) GetBlogByID(ctx context.Context, blogID string) (*entity.Blog, error) {
	var blog entity.Blog
	filter := bson.M{"_id": blogID, "is_deleted": false}

	err := r.collection.FindOne(ctx, filter).Decode(&blog)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("blog with id '%s' not found or has been deleted: %w", blogID, err)
		}
		return nil, fmt.Errorf("failed to retrieve blog post: %w", err)
	}

	return &blog, nil
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
	sortField := "created_at"
	if opts.SortBy != "" {
		sortField = opts.SortBy
	}

	// Use $facet for efficient counting and pagination in a single query
	facetStage := bson.M{
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
	}

	pipeline := []bson.M{
		{"$match": baseMatch},
		facetStage,
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
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

// UpdateBlog updates an existing blog post with new data.
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

	searchConditions := []bson.M{
		{"title": bson.M{"$regex": query, "$options": "i"}},
		{"authorDetails.username": bson.M{"$regex": query, "$options": "i"}},
		{"authorDetails.first_name": bson.M{"$regex": query, "$options": "i"}},
		{"authorDetails.last_name": bson.M{"$regex": query, "$options": "i"}},
	}

	if _, err := uuid.Parse(query); err == nil {
		searchConditions = append(searchConditions, bson.M{"author_id": query})
	}

	sortField := opts.SortBy
	switch sortField {
	case "":
		sortField = "created_at"
	case "username", "first_name", "last_name":
		sortField = "authorDetails." + sortField
	}

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
		{"$facet": bson.M{
			"totalCount": bson.A{
				bson.M{"$count": "total"},
			},
			"blogs": bson.A{
				bson.M{"$sort": bson.M{sortField: getSortOrder(opts.SortOrder)}},
				bson.M{"$skip": int64((opts.Page - 1) * opts.PageSize)},
				bson.M{"$limit": int64(opts.PageSize)},
			},
		}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
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
	filter := bson.M{"_id": blogID, "is_deleted": false, "like_count": bson.M{"$gt": 0}}
	update := bson.M{"$inc": bson.M{"like_count": -1}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to decrement like count: %w", err)
	}
	if res.ModifiedCount == 0 {
		return errors.New("blog post not found or like count is already zero")
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
// func (r *BlogRepository) DecrementDislikeCount(ctx context.Context, blogID string) error {
// 	filter := bson.M{"_id": blogID, "is_deleted": false}
// 	update := bson.M{"$inc": bson.M{"dislike_count": -1}}

// 	res, err := r.collection.UpdateOne(ctx, filter, update)
// 	if err != nil {
// 		return fmt.Errorf("failed to decrement dislike count: %w", err)
// 	}
// 	if res.ModifiedCount == 0 {
// 		return errors.New("blog post not found")
// 	}

// 	return nil
// }

// IncrementCommentCount increments the comment count of a specific blog post.
// func (r *BlogRepository) IncrementCommentCount(ctx context.Context, blogID string) error {
// 	filter := bson.M{"_id": blogID, "is_deleted": false}
// 	update := bson.M{"$inc": bson.M{"comment_count": 1}}

// 	res, err := r.collection.UpdateOne(ctx, filter, update)
// 	if err != nil {
// 		return fmt.Errorf("failed to increment comment count: %w", err)
// 	}
// 	if res.ModifiedCount == 0 {
// 		return errors.New("blog post not found")
// 	}

// 	return nil
// }

// DecrementCommentCount decrements the comment count of a specific blog post.
// func (r *BlogRepository) DecrementCommentCount(ctx context.Context, blogID string) error {
// 	filter := bson.M{"_id": blogID, "is_deleted": false}
// 	update := bson.M{"$inc": bson.M{"comment_count": -1}}

// 	res, err := r.collection.UpdateOne(ctx, filter, update)
// 	if err != nil {
// 		return fmt.Errorf("failed to decrement comment count: %w", err)
// 	}
// 	if res.ModifiedCount == 0 {
// 		return errors.New("blog post not found")
// 	}

// 	return nil
// }

// GetBlogCounts returns the current counts for a blog post.
func (r *BlogRepository) GetBlogCounts(ctx context.Context, blogID string) (viewCount, likeCount, dislikeCount, commentCount int, err error) {
	var blog entity.Blog
	filter := bson.M{"_id": blogID, "is_deleted": false}
	projection := bson.M{
		"view_count":    1,
		"like_count":    1,
		"dislike_count": 1,
		"comment_count": 1,
	}
	opts := options.FindOne().SetProjection(projection)

	err = r.collection.FindOne(ctx, filter, opts).Decode(&blog)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, 0, 0, 0, errors.New("blog not found")
		}
		return 0, 0, 0, 0, fmt.Errorf("failed to get blog counts: %w", err)
	}
	return blog.ViewCount, blog.LikeCount, blog.DislikeCount, blog.CommentCount, nil
}

// AddTagsToBlog associates one or more tags with a blog post.
func (r *BlogRepository) AddTagsToBlog(ctx context.Context, blogID string, tagIDs []string) error {
	if len(tagIDs) == 0 {
		return nil
	}

	filter := bson.M{"_id": blogID, "is_deleted": false}
	update := bson.M{
		"$addToSet": bson.M{"tags": bson.M{"$each": tagIDs}}, // Use $addToSet to avoid duplicate tags
		"$set":      bson.M{"updated_at": time.Now()},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to add tags to blog: %w", err)
	}
	if res.MatchedCount == 0 {
		return errors.New("blog post not found or already deleted")
	}

	return nil
}

// RemoveTagsFromBlog disassociates one or more tags from a blog post.
func (r *BlogRepository) RemoveTagsFromBlog(ctx context.Context, blogID string, tagIDs []string) error {
	if len(tagIDs) == 0 {
		return nil
	}

	filter := bson.M{"_id": blogID, "is_deleted": false}
	update := bson.M{
		"$pull": bson.M{"tags": bson.M{"$in": tagIDs}}, // Use $pull to remove items from the array
		"$set":  bson.M{"updated_at": time.Now()},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to remove tags from blog: %w", err)
	}
	if res.ModifiedCount == 0 {
		// This can happen if the blog exists but the tags to be removed aren't in the array.
		// We can consider this a successful operation in most cases.
		// To be strict, we first check if the blog exists.
		count, _ := r.collection.CountDocuments(ctx, filter)
		if count == 0 {
			return errors.New("blog post not found or already deleted")
		}
	}

	return nil
}

// GetBlogsByTagID retrieves a list of blog posts associated with a specific tag ID, applying pagination and sorting options.
func (r *BlogRepository) GetBlogsByTagID(ctx context.Context, tagID string, opts *contract.BlogFilterOptions) ([]*entity.Blog, int64, error) {
	// We simply set the TagIDs filter and delegate to the main GetBlogs function.
	// This avoids code duplication and keeps filtering logic centralized.
	filterOpts := &contract.BlogFilterOptions{
		Page:      opts.Page,
		PageSize:  opts.PageSize,
		SortBy:    opts.SortBy,
		SortOrder: opts.SortOrder,
		TagIDs:    []string{tagID}, // Filter by the specific tag ID
	}

	return r.GetBlogs(ctx, filterOpts)
}

func (r *BlogRepository) GetBlogsByTagIDs(ctx context.Context, tagIDs []string, page int, pageSize int) ([]*entity.Blog, int64, error) {
	if len(tagIDs) == 0 {
		return []*entity.Blog{}, 0, nil
	}

	// We delegate to the main GetBlogs function, which is designed to handle this case efficiently.
	filterOpts := &contract.BlogFilterOptions{
		Page:      page,
		PageSize:  pageSize,
		TagIDs:    tagIDs,
		SortBy:    "created_at", // Default sort order
		SortOrder: "desc",
	}

	return r.GetBlogs(ctx, filterOpts)
}

// HasViewedRecently checks if a user (by user ID or IP address) has viewed a blog within the last 24 hours.
func (r *BlogRepository) HasViewedRecently(ctx context.Context, blogID, userID, ipAddress string) (bool, error) {
	filter := bson.M{
		"blog_id": blogID,
		"$or": []bson.M{
			{"ip_address": ipAddress},
		},
	}
	if userID != "" {
		filter["$or"] = append(filter["$or"].([]bson.M), bson.M{"user_id": userID})
	}

	count, err := r.blogViewsCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("failed to check for recent blog view: %w", err)
	}
	return count > 0, nil
}

// RecordView records a user's view of a blog, including IP address and user agent.
func (r *BlogRepository) RecordView(ctx context.Context, blogID, userID, ipAddress, userAgent string) error {
	view := entity.BlogView{
		BlogID:    blogID,
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		ViewedAt:  time.Now(),
	}
	_, err := r.blogViewsCollection.InsertOne(ctx, view)
	if err != nil {
		return fmt.Errorf("failed to record blog view: %w", err)
	}
	return nil
}

// GetRecentViewsByIP retrieves recent views from a specific IP address.
func (r *BlogRepository) GetRecentViewsByIP(ctx context.Context, ipAddress string, since time.Time) ([]entity.BlogView, error) {
	filter := bson.M{
		"ip_address": ipAddress,
		"viewed_at":  bson.M{"$gte": since},
	}

	cursor, err := r.blogViewsCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve recent views by IP: %w", err)
	}
	defer cursor.Close(ctx)

	var views []entity.BlogView
	if err = cursor.All(ctx, &views); err != nil {
		return nil, fmt.Errorf("failed to decode recent views: %w", err)
	}

	return views, nil
}

// GetRecentViewsByUser retrieves all views from a specific user ID within a given time frame.
func (r *BlogRepository) GetRecentViewsByUser(ctx context.Context, userID string, since time.Time) ([]entity.BlogView, error) {
	if userID == "" {
		return []entity.BlogView{}, nil
	}

	filter := bson.M{
		"user_id":   userID,
		"viewed_at": bson.M{"$gte": since},
	}

	cursor, err := r.blogViewsCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve recent views by user: %w", err)
	}
	defer cursor.Close(ctx)

	var views []entity.BlogView
	if err = cursor.All(ctx, &views); err != nil {
		return nil, fmt.Errorf("failed to decode recent views: %w", err)
	}

	return views, nil
}
