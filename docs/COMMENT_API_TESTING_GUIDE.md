# Comment API Testing Guide

## Postman Collection Examples

### 1. Authentication Required Endpoints

#### Create Comment on Blog
```http
POST /api/v1/blogs/{{blogId}}/comments
Authorization: Bearer {{jwt_token}}
Content-Type: application/json

{
  "content": "This is a great blog post! Thanks for sharing."
}
```

#### Create Reply to Comment
```http
POST /api/v1/comments/{{commentId}}/replies
Authorization: Bearer {{jwt_token}}
Content-Type: application/json

{
  "content": "I completely agree with your point!"
}
```

#### Update Comment (within 15 min window)
```http
PUT /api/v1/comments/{{commentId}}
Authorization: Bearer {{jwt_token}}
Content-Type: application/json

{
  "content": "Updated: This is an even better blog post!"
}
```

#### Delete Comment
```http
DELETE /api/v1/comments/{{commentId}}
Authorization: Bearer {{jwt_token}}
```

#### Like/Unlike Comment
```http
POST /api/v1/comments/{{commentId}}/like
Authorization: Bearer {{jwt_token}}
```

### 2. Public Endpoints (No Authentication)

#### Get Blog Comments
```http
GET /api/v1/blogs/{{blogId}}/comments?page=1&limit=20&sort=newest
```

#### Get Specific Comment
```http
GET /api/v1/comments/{{commentId}}
```

#### Get Comment Thread
```http
GET /api/v1/comments/{{commentId}}/thread?max_depth=5
```

#### Get Comment Replies
```http
GET /api/v1/comments/{{commentId}}/replies?page=1&limit=10
```

#### Get Comment Statistics
```http
GET /api/v1/comments/{{commentId}}/stats
```

#### Get Comment Depth
```http
GET /api/v1/comments/{{commentId}}/depth
```

#### Get Blog Comment Count
```http
GET /api/v1/blogs/{{blogId}}/comments/count
```

#### Search Comments
```http
GET /api/v1/comments/search?q=great&blog_id={{blogId}}&status=approved
```

#### Get User's Comments
```http
GET /api/v1/users/{{userId}}/comments?page=1&limit=20
```

### 3. User-Specific Endpoints

#### Get Current User's Comments
```http
GET /api/v1/users/me/comments
Authorization: Bearer {{jwt_token}}
```

### 4. Moderation & Reporting

#### Report Comment
```http
POST /api/v1/comments/{{commentId}}/report
Authorization: Bearer {{jwt_token}}
Content-Type: application/json

{
  "reason": "spam",
  "details": "This comment contains spam content"
}
```

### 5. Admin Endpoints

#### Bulk Delete Comments
```http
DELETE /api/v1/admin/comments/bulk
Authorization: Bearer {{admin_jwt_token}}
Content-Type: application/json

{
  "comment_ids": ["uuid1", "uuid2", "uuid3"],
  "reason": "Violates community guidelines"
}
```

#### Update Comment Status
```http
PUT /api/v1/admin/comments/{{commentId}}/status
Authorization: Bearer {{admin_jwt_token}}
Content-Type: application/json

{
  "status": "hidden"
}
```

#### Get Comment Reports
```http
GET /api/v1/admin/comments/reports?page=1&page_size=20
Authorization: Bearer {{admin_jwt_token}}
```

## Expected Response Formats

### Comment Response
```json
{
  "comment": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "blog_id": "123e4567-e89b-12d3-a456-426614174001",
    "parent_id": null,
    "author_id": "123e4567-e89b-12d3-a456-426614174002",
    "author_name": "John Doe",
    "content": "This is a great blog post!",
    "status": "approved",
    "like_count": 5,
    "reply_count": 2,
    "is_liked": false,
    "created_at": "2025-08-06T10:30:00Z",
    "updated_at": "2025-08-06T10:30:00Z"
  }
}
```

### Comments List Response
```json
{
  "comments": [...],
  "pagination": {
    "current_page": 1,
    "page_size": 20,
    "total_items": 150,
    "total_pages": 8,
    "has_next": true,
    "has_previous": false
  }
}
```

### Error Response
```json
{
  "error": "Content is required",
  "details": "Comment content cannot be empty"
}
```

## Rate Limiting

- Comment creation/replies: 30 requests per minute
- Other operations: Standard rate limits apply
- Rate limit headers included in responses

## Testing Scenarios

### 1. Content Validation Tests
- Empty content → 400 Bad Request
- Content > 1000 chars → 400 Bad Request
- HTML content → Sanitized
- Profanity → 400 Bad Request

### 2. Authentication Tests
- No token → 401 Unauthorized
- Invalid token → 401 Unauthorized
- Expired token → 401 Unauthorized

### 3. Authorization Tests
- Edit other user's comment → 403 Forbidden
- Admin accessing user comment → 200 OK
- User accessing admin endpoint → 403 Forbidden

### 4. Rate Limiting Tests
- 31 requests in 1 minute → 429 Too Many Requests
- Request after rate limit reset → 200 OK

### 5. Business Logic Tests
- Reply to non-existent comment → 404 Not Found
- Edit after 15 minutes → 403 Forbidden
- Create comment on non-existent blog → 404 Not Found

## Environment Variables

```env
MONGODB_URI=mongodb://localhost:27017
MONGODB_DB_NAME=blog_db
JWT_SECRET=your-secret-key
PORT=8080
```

## Running the API

```bash
# Build the application
go build ./cmd/api

# Run with environment variables
./api

# Or run directly
go run ./cmd/api/main.go
```

The API will start on port 8080 (or PORT env variable) and be ready for testing!
