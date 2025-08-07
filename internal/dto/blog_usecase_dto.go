package dto

import (
    "time"
)

// Request DTOs
type CreateCommentRequest struct {
    Content  string     `json:"content" validate:"required,min=1,max=1000"`
    ParentID *string `json:"parent_id"`
    TargetID *string `json:"target_id"`
}

type UpdateCommentRequest struct {
    Content string `json:"content" validate:"required,min=1,max=1000"`
}

type UpdateCommentStatusRequest struct {
    Status string `json:"status" validate:"required,oneof=approved pending hidden flagged"`
}

type ReportCommentRequest struct {
    Reason  string `json:"reason" validate:"required,oneof=spam harassment inappropriate offensive"`
    Details string `json:"details" validate:"max=500"`
}

// Response DTOs
type CommentResponse struct {
    ID         string  `json:"id"`
    BlogID     string  `json:"blog_id"`
    ParentID   *string `json:"parent_id"`
    TargetID   *string `json:"target_id"`
    AuthorID   string  `json:"author_id"`
    AuthorName string `json:"author_name"`
    Content    string     `json:"content"`
    Status     string     `json:"status"`
    LikeCount  int        `json:"like_count"`
    IsLiked    bool       `json:"is_liked"`
    CreatedAt  time.Time  `json:"created_at"`
    UpdatedAt  time.Time  `json:"updated_at"`
    ReplyCount int        `json:"reply_count"`
}

type CommentThreadResponse struct {
    Comment *CommentResponse         `json:"comment"`
    Replies []*CommentThreadResponse `json:"replies"`
    Depth   int                      `json:"depth"`
}

type CommentsResponse struct {
    Comments   []*CommentResponse `json:"comments"`
    Pagination PaginationMeta     `json:"pagination"`
}

type PaginationMeta struct {
    CurrentPage  int   `json:"current_page"`
    PageSize     int   `json:"page_size"`
    TotalItems   int64 `json:"total_items"`
    TotalPages   int   `json:"total_pages"`
    HasNext      bool  `json:"has_next"`
    HasPrevious  bool  `json:"has_previous"`
}

type CommentReportResponse struct {
    ID         string  `json:"id"`
    CommentID  string  `json:"comment_id"`
    ReporterID string  `json:"reporter_id"`
    Reason     string     `json:"reason"`
    Details    string     `json:"details"`
    Status     string     `json:"status"`
    CreatedAt  time.Time  `json:"created_at"`
    ReviewedAt *time.Time `json:"reviewed_at"`
    ReviewedBy *string `json:"reviewed_by"`
}

type ReportsResponse struct {
    Reports    []*CommentReportResponse `json:"reports"`
    Pagination PaginationMeta           `json:"pagination"`
}