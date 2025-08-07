package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/contract"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
)

// SortOrder defines sorting direction for list queries
type SortOrder string

const (
   SortOrderASC  SortOrder = "asc"
   SortOrderDESC SortOrder = "desc"
)

// IBlogUseCase defines blog-related business logic
type IBlogUseCase interface {
	CreateBlog(ctx context.Context, title, content string, authorID string, slug string, status BlogStatus, featuredImageID *string) (*entity.Blog, error)
	GetBlogs(ctx context.Context, page, pageSize int, sortBy string, sortOrder SortOrder, dateFrom *time.Time, dateTo *time.Time) (blogs []entity.Blog, totalCount int, currentPage int, totalPages int, err error)
	GetBlogDetail(cnt context.Context, slug string) (blog entity.Blog, err error)
	UpdateBlog(ctx context.Context, blogID, authorID string, title *string, content *string, status *BlogStatus, featuredImageID *string) (*entity.Blog, error)
	DeleteBlog(ctx context.Context, blogID, userID string, isAdmin bool) (bool, error)
   SearchAndFilterBlogs(ctx context.Context, query string, tags []string, dateFrom *time.Time, dateTo *time.Time, minViews *int, maxViews *int, minLikes *int, maxLikes *int, authorID *string, page int, pageSize int) ([]entity.Blog, int, int, int, error)
   TrackBlogView(ctx context.Context, blogID, userID, ipAddress, userAgent string) error
   GetPopularBlogs(ctx context.Context, page, pageSize int) ([]entity.Blog, int, int, int, error)
}
// BlogStatus defines the state of a blog post
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

	// If slug is not provided, generate it from the title
	if slug == "" {
		slug = strings.ReplaceAll(strings.ToLower(title), " ", "-")
	}

	blog := &entity.Blog{
		ID:              uc.uuidgen.NewUUID(),
		Title:           title,
		Content:         content,
		AuthorID:        authorID,
		Slug:            slug + "-" + uc.uuidgen.NewUUID(), // A UUID is always appended to ensure the final slug is unique
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


// GetBlogDetail retrieves a blog by its slug
func (uc *BlogUseCaseImpl) GetBlogDetail(ctx context.Context, slug string) (entity.Blog, error) {
   if slug == "" {
	   return entity.Blog{}, errors.New("slug is required")
   }
   // Find blog by slug (case-insensitive)
   filterOptions := &contract.BlogFilterOptions{
	   Page:      1,
	   PageSize:  1,
   }
   blogs, _, err := uc.blogRepo.GetBlogs(ctx, filterOptions)
   if err != nil {
	   uc.logger.Errorf("failed to get blogs for detail: %v", err)
	   return entity.Blog{}, fmt.Errorf("failed to get blog: %w", err)
   }
   for _, b := range blogs {
	   if strings.EqualFold(b.Slug, slug) && !b.IsDeleted {
		   return *b, nil
	   }
   }
   return entity.Blog{}, errors.New("blog not found")
}

// UpdateBlog updates an existing blog post
func (uc *BlogUseCaseImpl) UpdateBlog(ctx context.Context, blogID, authorID string, title *string, content *string, status *BlogStatus, featuredImageID *string) (*entity.Blog, error) {
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
		// Generate a new slug from the new title
		newSlug := strings.ReplaceAll(strings.ToLower(*title), " ", "-")
		updates["slug"] = newSlug + "-" + uc.uuidgen.NewUUID()
	}
	if content != nil {
		updates["content"] = *content
	}

	if status != nil {
		updates["status"] = entity.BlogStatus(*status)
		if *status == BlogStatusPublished && blog.PublishedAt == nil {
			now := time.Now()
			updates["published_at"] = &now
		}
	}

	if featuredImageID != nil {
		updates["featured_image_id"] = *featuredImageID
	}

	if len(updates) > 0 {
		updates["updated_at"] = time.Now()
		if err := uc.blogRepo.UpdateBlog(ctx, blogID, updates); err != nil {
			uc.logger.Errorf("failed to update blog: %v", err)
			return nil, fmt.Errorf("failed to update blog: %w", err)
		}
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

// TrackBlogView tracks a view on a blog post, ensuring it's authentic by checking user ID, IP address, and User-Agent.

// isBot returns true if the User-Agent string matches common bot patterns.
func isBot(userAgent string) bool {
	ua := strings.ToLower(userAgent)
	botSignatures := []string{"bot", "spider", "crawl", "slurp", "curl", "wget", "python-requests", "httpclient", "feedfetcher", "mediapartners-google"}
	for _, sig := range botSignatures {
		if strings.Contains(ua, sig) {
			return true
		}
	}
	return false
}
func (uc *BlogUseCaseImpl) TrackBlogView(ctx context.Context, blogID, userID, ipAddress, userAgent string) error {
	if blogID == "" {
		return errors.New("blog ID is required")
	}

	// For a view to be considered unique, either the userID (if logged in) or the IP address must be provided.
	if userID == "" && ipAddress == "" {
		return errors.New("unable to track view without user ID or IP address")
	}

	// 1. Basic Bot Detection
	if isBot(userAgent) {
		uc.logger.Infof("Bot detected, view not counted for blog %s. User-Agent: %s", blogID, userAgent)
		return nil
	}

	// 2. Check for recent view from this user/IP for this specific blog post
	hasViewed, err := uc.blogRepo.HasViewedRecently(ctx, blogID, userID, ipAddress)
	if err != nil {
		uc.logger.Errorf("failed to check for recent blog view: %v", err)
		return fmt.Errorf("failed to check for recent blog view: %w", err)
	}
	if hasViewed {
		return nil // Already viewed this post recently
	}

	// 3. Advanced Velocity & Rotation Checks
	// Define time windows for checks
	shortWindow := time.Now().Add(-5 * time.Minute)  // for rapid-fire views - 5 minutes
	mediumWindow := time.Now().Add(-60 * time.Minute) // for IP rotation     - 60 minutes

	// Fetch recent activity
	ipViews, err := uc.blogRepo.GetRecentViewsByIP(ctx, ipAddress, shortWindow)
	if err != nil {
		return fmt.Errorf("failed to get recent views by IP: %w", err)
	}
	userViews, err := uc.blogRepo.GetRecentViewsByUser(ctx, userID, mediumWindow)
	if err != nil {
		return fmt.Errorf("failed to get recent views by user: %w", err)
	}

	// IP Velocity Check: Has this IP viewed too many different blogs in the last 5 minutes?
	const maxIpVelocity = 10 // Max 10 views from one IP in 5 mins
	if len(ipViews) > maxIpVelocity {
		uc.logger.Warningf("High IP velocity detected for %s. Views: %d", ipAddress, len(ipViews))
		return errors.New("suspicious activity detected: high view velocity")
	}

	// User-IP Rotation Check: Has this user account used too many IPs in the last hour?
	if userID != "" {
		const maxUserIPs = 5 // Max 5 different IPs for one user in 1 hour
		ipSet := make(map[string]struct{})
		for _, view := range userViews {
			ipSet[view.IPAddress] = struct{}{}
		}
		if len(ipSet) > maxUserIPs {
			uc.logger.Warningf("High IP rotation detected for user %s. IPs used: %d", userID, len(ipSet))
			return errors.New("suspicious activity detected: high IP rotation")
		}
	}

	// If all checks pass, increment the view count and record the view
	if err := uc.blogRepo.IncrementViewCount(ctx, blogID); err != nil {
		uc.logger.Errorf("failed to increment view count: %v", err)
		return fmt.Errorf("failed to increment view count: %w", err)
	}

	if err := uc.blogRepo.RecordView(ctx, blogID, userID, ipAddress, userAgent); err != nil {
		uc.logger.Errorf("failed to record user view: %v", err)
		return fmt.Errorf("failed to record user view: %w", err)
	}

	return nil
}

// GetPopularBlogs returns blogs sorted by view count (descending), paginated.
func (uc *BlogUseCaseImpl) GetPopularBlogs(ctx context.Context, page, pageSize int) ([]entity.Blog, int, int, int, error) {
    if page < 1 {
        page = 1
    }
    if pageSize < 1 {
        pageSize = 10
    }

    filterOptions := &contract.BlogFilterOptions{
        Page:      page,
        PageSize:  pageSize,
        SortBy:    "viewCount",
        SortOrder: "desc",
    }

    blogs, totalCount, err := uc.blogRepo.GetBlogs(ctx, filterOptions)
    if err != nil {
        uc.logger.Errorf("failed to get popular blogs: %v", err)
        return nil, 0, 0, 0, fmt.Errorf("failed to get popular blogs: %w", err)
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

// SearchAndFilterBlogs implements advanced search and filtering for blogs.
func (uc *BlogUseCaseImpl) SearchAndFilterBlogs(
    ctx context.Context,
    query string,
    tags []string,
    dateFrom *time.Time,
    dateTo *time.Time,
    minViews *int,
    maxViews *int,
    minLikes *int,
    maxLikes *int,
    authorID *string,
    page int,
    pageSize int,
) ([]entity.Blog, int, int, int, error) {
    filterOptions := &contract.BlogFilterOptions{
        Page:      page,
        PageSize:  pageSize,
        DateFrom:  dateFrom,
        DateTo:    dateTo,
        MinViews:  minViews,
        MaxViews:  maxViews,
        MinLikes:  minLikes,
        MaxLikes:  maxLikes,
        AuthorID:  authorID,
        TagIDs:    tags,
    }
    var blogs []*entity.Blog
    var totalCount int64
    var err error
    if query != "" {
        blogs, totalCount, err = uc.blogRepo.SearchBlogs(ctx, query, filterOptions)
    } else {
        blogs, totalCount, err = uc.blogRepo.GetBlogs(ctx, filterOptions)
    }
    if err != nil {
        uc.logger.Errorf("failed to search/filter blogs: %v", err)
        return nil, 0, 0, 0, fmt.Errorf("failed to search/filter blogs: %w", err)
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