# G6 Blog API Endpoints (Updated)

This document lists all available API endpoints for the G6 Blog Starter Project, including authentication, user management, blog management, and interactions. It is based on the Postman collection and the current router implementation.

---

## Authentication

- **POST** `/api/v1/auth/register` — Register a new user
- **POST** `/api/v1/auth/login` — Login and receive tokens
- **POST** `/api/v1/auth/verify-email` — Verify email with token
- **POST** `/api/v1/auth/forgot-password` — Request password reset
- **POST** `/api/v1/auth/reset-password` — Reset password with token
- **POST** `/api/v1/auth/refresh-token` — Refresh access token
- **POST** `/api/v1/logout` — Logout and invalidate refresh token
- **GET** `/auth/google/login` — Google OAuth2 login

## User Management

- **GET** `/api/v1/users/:id` — Get user by ID (auth required)
- **GET** `/api/v1/me` — Get current user (auth required)
- **PUT** `/api/v1/me` — Update current user (auth required)

## Blog Management

- **GET** `/api/v1/blogs` — List blogs (supports pagination, sorting, filtering)
- **GET** `/api/v1/blogs/search` — Search and filter blogs (query, tags, date, views, likes, author, pagination)
- **GET** `/api/v1/blogs/popular` — Get popular blogs (sorted by view count)
- **GET** `/api/v1/blogs/:slug` — Get blog details by slug
- **POST** `/api/v1/blogs` — Create a new blog (auth required)
- **PUT** `/api/v1/blogs/:blogID` — Update a blog (auth required)
- **DELETE** `/api/v1/blogs/:blogID` — Delete a blog (auth required)

## Blog Interactions

- **POST** `/api/v1/blogs/:blogID/like` — Like a blog (auth required)
- **DELETE** `/api/v1/blogs/:blogID/like` — Unlike a blog (auth required)
- **POST** `/api/v1/blogs/:blogID/view` — Track a blog view (auth required)

---

### Notes

- All endpoints under `/api/v1` are versioned.
- Endpoints marked as "auth required" require a valid JWT access token in the `Authorization` header.
- Pagination, sorting, and filtering parameters are supported on list endpoints.
- For detailed request/response examples, see the Postman collection: `api/postman/g6-blog.postman_collection.json`.
