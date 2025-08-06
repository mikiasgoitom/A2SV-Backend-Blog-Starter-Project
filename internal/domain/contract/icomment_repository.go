package contract

import (
    "context"
    "github.com/google/uuid"
    "github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/entity"
)

type Pagination struct {
    Page     int `json:"page"`
    PageSize int `json:"page_size"`
}

type PaginationMeta struct {
    CurrentPage  int   `json:"current_page"`
    PageSize     int   `json:"page_size"`
    TotalItems   int64 `json:"total_items"`
    TotalPages   int   `json:"total_pages"`
    HasNext      bool  `json:"has_next"`
    HasPrevious  bool  `json:"has_previous"`
}

type ICommentRepository interface {
    // Core CRUD operations
    Create(ctx context.Context, comment *entity.Comment) error
    GetByID(ctx context.Context, id uuid.UUID) (*entity.Comment, error)
    Update(ctx context.Context, comment *entity.Comment) error
    Delete(ctx context.Context, id uuid.UUID) error

    // Listing operations
    GetTopLevelComments(ctx context.Context, blogID uuid.UUID, pagination Pagination) ([]*entity.Comment, int64, error)
    GetCommentThread(ctx context.Context, parentID uuid.UUID) (*entity.CommentThread, error)
    GetCommentsByUser(ctx context.Context, userID uuid.UUID, pagination Pagination) ([]*entity.Comment, int64, error)

    // Status and moderation
    UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
    GetCommentCount(ctx context.Context, blogID uuid.UUID) (int64, error)

    // Like system
    LikeComment(ctx context.Context, commentID, userID uuid.UUID) error
    UnlikeComment(ctx context.Context, commentID, userID uuid.UUID) error
    IsCommentLikedByUser(ctx context.Context, commentID, userID uuid.UUID) (bool, error)
    GetCommentLikeCount(ctx context.Context, commentID uuid.UUID) (int64, error)

    // Reporting system
    ReportComment(ctx context.Context, report *entity.CommentReport) error
    GetCommentReports(ctx context.Context, pagination Pagination) ([]*entity.CommentReport, int64, error)
    UpdateReportStatus(ctx context.Context, reportID uuid.UUID, status string, reviewerID uuid.UUID) error
}
