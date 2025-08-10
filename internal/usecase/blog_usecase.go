package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/contract"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
	usecasecontract "github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/usecase/contract"
)

// BlogUseCaseImpl implements the BlogUseCase interface
type BlogUseCaseImpl struct {
	blogRepo contract.IBlogRepository
	uuidgen  contract.IUUIDGenerator
	logger   usecasecontract.IAppLogger
	aiUC     usecasecontract.IAIUseCase
}

// NewBlogUseCase creates a new instance of BlogUseCase
func NewBlogUseCase(blogRepo contract.IBlogRepository, uuidgenrator contract.IUUIDGenerator, logger usecasecontract.IAppLogger, aiUC usecasecontract.IAIUseCase) *BlogUseCaseImpl {
	return &BlogUseCaseImpl{
		blogRepo: blogRepo,
		logger:   logger,
		uuidgen:  uuidgenrator,
		aiUC:     aiUC,
	}
}

// check if UserUseCase implements the IUserUseCase
var _ usecasecontract.IBlogUseCase = (*BlogUseCaseImpl)(nil)

// CreateBlog creates a new blog post
func (uc *BlogUseCaseImpl) CreateBlog(ctx context.Context, title, content string, authorID string, slug string, status usecasecontract.BlogStatus, featuredImageID *string, tags []string) (*entity.Blog, error) {
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
		LikeCount:       0,
		DislikeCount:    0,
		CommentCount:    0,
		Popularity:      calculatePopularity(0, 0, 0, 0),
		FeaturedImageID: featuredImageID,
		IsDeleted:       false,
	}

	if status == usecasecontract.BlogStatusPublished {
		now := time.Now()
		blog.PublishedAt = &now
	}
	// Check for profanity in the content
	feedback, err := uc.aiUC.CensorAndCheckBlog(ctx, content)
	if err != nil {
		return nil, fmt.Errorf("failed to check content: %w", err)
	}
	if feedback == "no" {
		return nil, errors.New("content contains inappropriate material")
	}

	if err := uc.blogRepo.CreateBlog(ctx, blog); err != nil {
		uc.logger.Errorf("failed to create blog: %v", err)
		return nil, fmt.Errorf("failed to create blog: %w", err)
	}
	// Add tags to blog if provided
	if len(tags) > 0 {
		err := uc.blogRepo.AddTagsToBlog(ctx, blog.Slug, tags)
		if err != nil {
			uc.logger.Errorf("Failed to add tags to blog: %v", err)
			// Not returning error here to allow blog creation to succeed even if tag association fails
		}
	}

	return blog, nil
}

// GetBlogs retrieves paginated list of blogs
func (uc *BlogUseCaseImpl) GetBlogs(ctx context.Context, page, pageSize int, sortBy string, sortOrder usecasecontract.SortOrder, dateFrom *time.Time, dateTo *time.Time) ([]entity.Blog, int, int, int, error) {
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

	// Only return published or archived blogs (not drafts)
	blogs, totalCount, err := uc.blogRepo.GetBlogs(ctx, filterOptions)
	if err != nil {
		uc.logger.Errorf("failed to get blogs: %v", err)
		return nil, 0, 0, 0, fmt.Errorf("failed to get blogs: %w", err)
	}

	var filteredBlogs []entity.Blog
	for _, blog := range blogs {
		if blog.Status == entity.BlogStatusPublished || blog.Status == entity.BlogStatusArchived {
			filteredBlogs = append(filteredBlogs, *blog)
		}
	}

	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize != 0 {
		totalPages++
	}

	return filteredBlogs, int(totalCount), page, totalPages, nil
}

// GetBlogDetail retrieves a blog by its slug
func (uc *BlogUseCaseImpl) GetBlogDetail(ctx context.Context, slug string) (entity.Blog, error) {
	if slug == "" {
		return entity.Blog{}, errors.New("slug is required")
	}
	blog, err := uc.blogRepo.GetBlogBySlug(ctx, slug)
	if err != nil {
		uc.logger.Errorf("failed to get blog by slug: %v", err)
		return entity.Blog{}, fmt.Errorf("failed to get blog: %w", err)
	}
	if blog == nil || blog.IsDeleted {
		return entity.Blog{}, errors.New("blog not found")
	}
	// Only allow published or archived blogs to be fetched by slug
	if blog.Status != entity.BlogStatusPublished && blog.Status != entity.BlogStatusArchived {
		return entity.Blog{}, errors.New("blog not found")
	}
	return *blog, nil
}

// UpdateBlog updates an existing blog post
func (uc *BlogUseCaseImpl) UpdateBlog(ctx context.Context, blogID, authorID string, title *string, content *string, status *usecasecontract.BlogStatus, featuredImageID *string) (*entity.Blog, error) {
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
		// if content is edited check for profanity
		feedback, err := uc.aiUC.CensorAndCheckBlog(ctx, *content)
		if err != nil {
			return nil, fmt.Errorf("failed to check content: %w", err)
		}
		if feedback == "no" {
			return nil, errors.New("content contains inappropriate material")
		}
	}

	if status != nil {
		updates["status"] = entity.BlogStatus(*status)
		if *status == usecasecontract.BlogStatusPublished && blog.PublishedAt == nil {
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
	shortWindow := time.Now().Add(-5 * time.Minute)   // for rapid-fire views - 5 minutes
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

	// Update popularity after view
	if err := uc.UpdateBlogPopularity(ctx, blogID); err != nil {
		uc.logger.Errorf("failed to update blog popularity after view: %v", err)
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
		SortBy:    "popularity",
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
		Page:     page,
		PageSize: pageSize,
		DateFrom: dateFrom,
		DateTo:   dateTo,
		MinViews: minViews,
		MaxViews: maxViews,
		MinLikes: minLikes,
		MaxLikes: maxLikes,
		AuthorID: authorID,
		TagIDs:   tags,
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

// calculatePopularity computes the popularity score for a blog
func calculatePopularity(views, likes, dislikes, comments int) float64 {
	// You can tune these weights as needed
	const (
		viewWeight    = 1.0
		likeWeight    = 3.0
		dislikeWeight = -2.0
		commentWeight = 2.0
	)
	return float64(views)*viewWeight + float64(likes)*likeWeight + float64(dislikes)*dislikeWeight + float64(comments)*commentWeight
}

// UpdateBlogPopularity fetches counts and updates the popularity field in the DB
func (uc *BlogUseCaseImpl) UpdateBlogPopularity(ctx context.Context, blogID string) error {
	views, likes, dislikes, comments, err := uc.blogRepo.GetBlogCounts(ctx, blogID)
	if err != nil {
		return err
	}
	popularity := calculatePopularity(views, likes, dislikes, comments)
	updates := map[string]interface{}{"popularity": popularity}
	return uc.blogRepo.UpdateBlog(ctx, blogID, updates)
}
