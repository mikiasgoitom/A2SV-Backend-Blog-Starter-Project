package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/contract"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
	"time"
)

type BlogUseCase interface {
	CreateBlog(ctx context.Context, title, content string, authorID string, slug string, status BlogStatus, featuredImageID *string) (*entity.Blog, error)
	GetBlogs(ctx context.Context, page, pageSize int, sortBy string, sortOrder SortOrder, dateFrom *time.Time, dateTo *time.Time) (blogs []entity.Blog, totalCount int, currentPage int, totalPages int, err error)
	UpdateBlog(ctx context.Context, blogID, authorID string, title *string, content *string, slug *string, status *BlogStatus, publishedAt *time.Time, featuredImageID *string, isDeleted *bool) (*entity.Blog, error)
	DeleteBlog(ctx context.Context, blogID, userID string, isAdmin bool) (bool, error)
	SearchAndFilterBlogs(ctx context.Context, query string, searchBy string, tags []string, dateFrom *time.Time, dateTo *time.Time, minViews *int, maxViews *int, minLikes *int, maxLikes *int, authorID *string, page int, pageSize int) (blogs []entity.Blog, err error, totalCount int, currentPage int, totalPages int)
	TrackBlogPopularity(ctx context.Context, blogID, userID string, action BlogAction) (viewCount, likeCount, dislikeCount, commentCount int, err error)
	GetRecommendedBlogs(ctx context.Context, userID string, page, pageSize int) (blogs []entity.Blog, err error)
}
type SortOrder string

const (
	SortOrderASC  SortOrder = "asc"
	SortOrderDESC SortOrder = "desc"
)

type BlogAction string

const (
	BlogActionView    BlogAction = "view"
	BlogActionLike    BlogAction = "like"
	BlogActionDislike BlogAction = "dislike"
	BlogActionComment BlogAction = "comment"
)

type BlogStatus string

const (
	BlogStatusDraft     BlogStatus = "draft"
	BlogStatusPublished BlogStatus = "published"
	BlogStatusArchived  BlogStatus = "archived"
)

// BlogUseCaseImpl implements the BlogUseCase interface
type BlogUseCaseImpl struct {
	blogRepo contract.IBlogRepository
	uuidgen  contract.IUUIDGenerator
	logger   AppLogger
}

// NewBlogUseCase creates a new instance of BlogUseCase
func NewBlogUseCase(blogRepo contract.IBlogRepository, uuidgenrator contract.IUUIDGenerator, logger AppLogger) *BlogUseCaseImpl {
	return &BlogUseCaseImpl{
		blogRepo: blogRepo,
		logger:   logger,
		uuidgen:  uuidgenrator,
	}
}

// CreateBlog creates a new blog post
func (uc *BlogUseCaseImpl) CreateBlog(ctx context.Context, title, content string, authorID string, slug string, status BlogStatus, featuredImageID *string) (*entity.Blog, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}
	if content == "" {
		return nil, errors.New("content is required")
	}
	if authorID == "" {
		return nil, errors.New("author ID is required")
	}
	if slug == "" {
		return nil, errors.New("slug is required")
	}

	blog := &entity.Blog{
		ID:              uc.uuidgen.NewUUID(),
		Title:           title,
		Content:         content,
		AuthorID:        authorID,
		Slug:            slug,
		Status:          entity.BlogStatus(status),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		ViewCount:       0,
		FeaturedImageID: featuredImageID,
		IsDeleted:       false,
	}

	if status == BlogStatusPublished {
		now := time.Now()
		blog.PublishedAt = &now
	}

	if err := uc.blogRepo.CreateBlog(ctx, blog); err != nil {
		uc.logger.Errorf("failed to create blog: %v", err)
		return nil, fmt.Errorf("failed to create blog: %w", err)
	}

	return blog, nil
}

// GetBlogs retrieves paginated list of blogs
func (uc *BlogUseCaseImpl) GetBlogs(ctx context.Context, page, pageSize int, sortBy string, sortOrder SortOrder, dateFrom *time.Time, dateTo *time.Time) ([]entity.Blog, int, int, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	filterOptions := &contract.BlogFilterOptions{
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: string(sortOrder),
		DateFrom:  dateFrom,
		DateTo:    dateTo,
	}

	blogs, totalCount, err := uc.blogRepo.GetBlogs(ctx, filterOptions)
	if err != nil {
		uc.logger.Errorf("failed to get blogs: %v", err)
		return nil, 0, 0, 0, fmt.Errorf("failed to get blogs: %w", err)
	}

	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize != 0 {
		totalPages++
	}

	var blogEntities []entity.Blog
	for _, blog := range blogs {
		blogEntities = append(blogEntities, *blog)
	}

	return blogEntities, int(totalCount), page, totalPages, nil
}

// UpdateBlog updates an existing blog post
func (uc *BlogUseCaseImpl) UpdateBlog(ctx context.Context, blogID, authorID string, title *string, content *string, slug *string, status *BlogStatus, publishedAt *time.Time, featuredImageID *string, isDeleted *bool) (*entity.Blog, error) {
	if blogID == "" {
		return nil, errors.New("blog ID is required")
	}
	if authorID == "" {
		return nil, errors.New("author ID is required")
	}

	// Get existing blog
	blog, err := uc.blogRepo.GetBlogByID(ctx, blogID)
	if err != nil {
		uc.logger.Errorf("failed to get blog: %v", err)
		return nil, fmt.Errorf("failed to get blog: %w", err)
	}
	if blog == nil {
		return nil, errors.New("blog not found")
	}

	// Check if user is the author
	if blog.AuthorID != authorID {
		return nil, errors.New("unauthorized: only the author can update this blog")
	}

	updates := make(map[string]interface{})

	if title != nil {
		updates["title"] = *title
	}
	if content != nil {
		updates["content"] = *content
	}
	if slug != nil {
		updates["slug"] = *slug
	}
	if status != nil {
		updates["status"] = entity.BlogStatus(*status)
		if *status == BlogStatusPublished && blog.PublishedAt == nil {
			now := time.Now()
			updates["published_at"] = &now
		}
	}
	if publishedAt != nil {
		updates["published_at"] = *publishedAt
	}
	if featuredImageID != nil {
		updates["featured_image_id"] = *featuredImageID
	}
	if isDeleted != nil {
		updates["is_deleted"] = *isDeleted
	}

	updates["updated_at"] = time.Now()

	if err := uc.blogRepo.UpdateBlog(ctx, blogID, updates); err != nil {
		uc.logger.Errorf("failed to update blog: %v", err)
		return nil, fmt.Errorf("failed to update blog: %w", err)
	}

	// Return updated blog
	updatedBlog, err := uc.blogRepo.GetBlogByID(ctx, blogID)
	if err != nil {
		uc.logger.Errorf("failed to get updated blog: %v", err)
		return nil, fmt.Errorf("failed to get updated blog: %w", err)
	}

	return updatedBlog, nil
}

// DeleteBlog deletes a blog post
func (uc *BlogUseCaseImpl) DeleteBlog(ctx context.Context, blogID, userID string, isAdmin bool) (bool, error) {
	if blogID == "" {
		return false, errors.New("blog ID is required")
	}
	if userID == "" {
		return false, errors.New("user ID is required")
	}

	blog, err := uc.blogRepo.GetBlogByID(ctx, blogID)
	if err != nil {
		uc.logger.Errorf("failed to get blog: %v", err)
		return false, fmt.Errorf("failed to get blog: %w", err)
	}
	if blog == nil {
		return false, errors.New("blog not found")
	}

	// Check authorization
	if !isAdmin && blog.AuthorID != userID {
		return false, errors.New("unauthorized: only the author or admin can delete this blog")
	}

	if err := uc.blogRepo.DeleteBlog(ctx, blogID); err != nil {
		uc.logger.Errorf("failed to delete blog: %v", err)
		return false, fmt.Errorf("failed to delete blog: %w", err)
	}

	return true, nil
}

// TrackBlogPopularity tracks blog interactions
func (uc *BlogUseCaseImpl) TrackBlogPopularity(ctx context.Context, blogID, userID string, action BlogAction) (int, int, int, int, error) {
	if blogID == "" {
		return 0, 0, 0, 0, errors.New("blog ID is required")
	}

	blog, err := uc.blogRepo.GetBlogByID(ctx, blogID)
	if err != nil {
		uc.logger.Errorf("failed to get blog: %v", err)
		return 0, 0, 0, 0, fmt.Errorf("failed to get blog: %w", err)
	}
	if blog == nil {
		return 0, 0, 0, 0, errors.New("blog not found")
	}

	viewCount := blog.ViewCount
	likeCount := 0
	dislikeCount := 0
	commentCount := 0

	switch action {
	case BlogActionView:
		if err := uc.blogRepo.IncrementViewCount(ctx, blogID); err != nil {
			uc.logger.Errorf("failed to increment view count: %v", err)
			return 0, 0, 0, 0, fmt.Errorf("failed to increment view count: %w", err)
		}
		viewCount++
	case BlogActionLike:
		likeCount++
	case BlogActionDislike:
		dislikeCount++
	case BlogActionComment:
		commentCount++
	}

	return viewCount, likeCount, dislikeCount, commentCount, nil
}
func (uc *BlogUseCaseImpl) SearchAndFilterBlogs(ctx context.Context, query string, searchBy string, tags []string, dateFrom *time.Time, dateTo *time.Time, minViews *int, maxViews *int, minLikes *int, maxLikes *int, authorID *string, page int, pageSize int) ([]entity.Blog, error, int, int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	var blogs []*entity.Blog
	var totalCount int64
	var err error

	if len(tags) > 0 {
		// Fetch blogs by tag IDs
		blogs, totalCount, err = uc.blogRepo.GetBlogsByTagIDs(ctx, tags, page, pageSize)
		if err != nil {
			uc.logger.Errorf("failed to get blogs by tag IDs: %v", err)
			return nil, fmt.Errorf("failed to get blogs by tag IDs: %w", err), 0, 0, 0
		}
	} else {
		filterOptions := &contract.BlogFilterOptions{
			Page:      page,
			PageSize:  pageSize,
			SortBy:    "created_at",
			SortOrder: "desc",
			DateFrom:  dateFrom,
			DateTo:    dateTo,
			MinViews:  minViews,
			MaxViews:  maxViews,
			MinLikes:  minLikes,
			MaxLikes:  maxLikes,
			AuthorID:  authorID,
		}

		// Use SearchBlogs if query is provided, otherwise use GetBlogs
		if query != "" {
			blogs, totalCount, err = uc.blogRepo.SearchBlogs(ctx, query, filterOptions)
		} else {
			blogs, totalCount, err = uc.blogRepo.GetBlogs(ctx, filterOptions)
		}

		if err != nil {
			uc.logger.Errorf("failed to search and filter blogs: %v", err)
			return nil, fmt.Errorf("failed to search and filter blogs: %w", err), 0, 0, 0
		}
	}

	// Post-processing for additional filters
	var filteredBlogs []*entity.Blog
	for _, blog := range blogs {
		// Skip deleted blogs
		if blog.IsDeleted {
			continue
		}

		filteredBlogs = append(filteredBlogs, blog)
	}

	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize != 0 {
		totalPages++
	}

	var blogEntities []entity.Blog
	for _, blog := range filteredBlogs {
		blogEntities = append(blogEntities, *blog)
	}

	return blogEntities, nil, int(totalCount), page, totalPages
}

// GetRecommendedBlogs gets recommended blogs for a user
func (uc *BlogUseCaseImpl) GetRecommendedBlogs(ctx context.Context, userID string, page, pageSize int) ([]entity.Blog, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	filterOptions := &contract.BlogFilterOptions{
		Page:      page,
		PageSize:  pageSize,
		SortBy:    "view_count",
		SortOrder: "desc",
	}

	blogs, _, err := uc.blogRepo.GetBlogs(ctx, filterOptions)
	if err != nil {
		uc.logger.Errorf("failed to get recommended blogs: %v", err)
		return nil, fmt.Errorf("failed to get recommended blogs: %w", err)
	}

	var blogEntities []entity.Blog
	for _, blog := range blogs {
		blogEntities = append(blogEntities, *blog)
	}

	return blogEntities, nil
}
