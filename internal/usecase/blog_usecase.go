// package usecase

// import (
// 	"context"
// 	"time"

// 	"github.com/google/uuid"
// 	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
// )

// func CreateBlog(ctx context.Context, blog entity.Blog ) (*entity.Blog, error){

	
// }


// func GetBlogByID(ctx context.Context, blogID uuid.UUID) (*entity.Blog, error)
// func UpdateBlog(ctx context.Context, blogID, authorID uuid.UUID, title *string, content *string, slug *string, status *entity.BlogStatus, publishedAt *time.Time, featuredImageID *uuid.UUID, isDeleted *bool) (*entity.Blog, error)
// func TrackBlogPopularity(ctx context.Context, blogID, userID uuid.UUID, action BlogAction) (viewCount, likeCount, dislikeCount, commentCount int, err error)
// func DeleteBlog(ctx context.Context, blogID, userID uuid.UUID, isAdmin bool) (bool, error)
// func GetBlogs(ctx context.Context, page, pageSize int, sortBy string, sortOrder SortOrder, dateFrom *time.Time, dateTo *time.Time) (blogs []entity.Blog, totalCount int, currentPage int, totalPages int, err error)
// func SearchAndFilterBlogs(ctx context.Context, query string, page, pageSize int, searchBy string, tags []string, dateFrom *time.Time, dateTo *time.Time, minViews *int, minLikes *int, authorID *uuid.UUID) (blogs []entity.Blog, totalCount int, currentPage int, totalPages int, err error)