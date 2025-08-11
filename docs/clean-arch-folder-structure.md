```
g67-blog-project-g1/
├── api/
│   └── postman/
│       └── g6-blog.postman_collection.json
├── cmd/
│   └── api/
│       └── main.go
├── configs/
│   └── config.yml
├── internal/
│   ├── entity/
│   │   ├── user.go
│   │   ├── blog.go
│   │   └── token.go
│   ├── usecase/
│   │   ├── interfaces.go
│   │   ├── user_usecase.go
│   │   ├── blog_usecase.go
│   │   └── ai_usecase.go
│   ├── handler/
│   │   └── http/
│   │       ├── middleware/
│   │       │   ├── auth.go
│   │       │   ├── ratelimiter.go
│   │       │   └── logger.go
│   │       ├── router.go
│   │       ├── user_handler.go
│   │       ├── blog_handler.go
│   │       ├── request.go
│   │       └── response.go
│   └── infrastructure/
│       ├── config/
│       ├── database/
│       ├── hash/
│       ├── jwt/
│       ├── logger/
│       ├── validator/
│       ├── repository/
│       │   └── mongodb/
│       │       ├── user_repo.go
│       │       ├── blog_repo.go
│       │       └── token_repo.go
│       └── external_services/
│           ├── ai_adapter.go
│           └── mail_service.go
├── migrations/
│   └── 000001_create_initial_indexes.go
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```
