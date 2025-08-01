Here is a comprehensive guide detailing the folder structure, naming conventions, GitHub collaboration, and commit conventions for your G6 Blog Starter Project, presented in a format suitable for a PDF document.

-----

# G6 Blog Starter Project: Development Guidelines

## 1\. Project Overview

This document outlines the architectural structure, naming conventions, and collaboration workflows for the G6 Blog Starter Project, a backend API for a blog platform. The project aims to provide robust CRUD operations for blog posts, comprehensive user management with authentication and authorization (including roles), and AI integration for content suggestions.

The architecture adheres to Robert C. Martin's Clean Architecture principles, promoting separation of concerns, testability, and scalability, particularly for Go applications utilizing MongoDB.

## 2\. Go Project Folder Structure (Clean Architecture)

The following folder structure organizes your Go backend API by feature and layer, promoting separation of concerns, testability, and scalability. It's designed to keep the core business logic independent of external frameworks and databases.

``` markdown
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

### Directory Explanations:

  * **`api/`**:
      * **Purpose:** Stores API-related documentation and specifications.
      * **Content:** `postman/g6-blog.postman_collection.json` for API documentation, example requests, and responses. If using OpenAPI/Swagger, the spec file would also reside here.
  * **`cmd/`**:
      * **Purpose:** Contains the main entry point for your application.
      * **Content:** `api/main.go` will be the executable file. It's responsible for initializing and wiring together all application components (config, logger, database connections, repositories, use cases, HTTP server). This is where **Dependency Injection** occurs, assembling the application's dependencies at startup.
  * **`configs/`**:
      * **Purpose:** Holds application configuration files.
      * **Content:** `config.yml` (or `.json`, `.env`) for environment-specific settings (database connection strings, API keys, server ports, etc.).
  * **`internal/`**:
      * **Purpose:** This directory holds all the private code for your application. The Go toolchain prevents other projects from importing packages from an `internal` directory, ensuring your application's core logic remains encapsulated.
      * **`internal/entity/` (The Core 💎)**:
          * **Purpose:** Contains the core domain models and business objects. These are simple Go structs representing fundamental concepts like `User`, `Blog`, or `Token`. They define the essential data structures for your application.
          * **Clean Architecture Layer:** **Entities**.
          * **Key Rule:** This layer is pure and has zero dependencies on any other layer in your project.
      * **`internal/usecase/` (The Brains 🧠)**:
          * **Purpose:** Contains the application-specific business logic. Each file coordinates a specific use case, like `user_usecase.go` for user registration or `blog_usecase.go` for blog creation. They define *what* the application does.
          * **Clean Architecture Layer:** **Use Cases**.
          * **`interfaces.go`**: This is a crucial file. It defines the **interfaces (contracts)** for your repositories (e.g., `UserRepository`, `BlogRepository`) and external services (`AIService`). The use cases depend on these interfaces, not on concrete implementations, adhering to the Dependency Inversion Principle.
      * **`internal/repository/` (The Data Layer 🗄️)**:
          * **Purpose:** Implements the repository interfaces defined in the `usecase` layer. This is where the actual database interactions and logic live.
          * **`internal/repository/mongodb/`**: Specifically for MongoDB implementations. Files like `user_repo.go`, `blog_repo.go`, `token_repo.go` will contain direct MongoDB driver code (without an ODM/ORM), building BSON documents for queries and updates, and decoding results manually.
          * **Clean Architecture Layer:** **Interface Adapters (Gateway Implementations)**. This layer isolates all MongoDB-specific code.
      * **`internal/handler/` (The Entryway 🚪)**:
          * **Purpose:** Handles incoming requests (e.g., from an HTTP API). It's responsible for parsing requests, calling the appropriate use case, and formatting the response.
          * **`internal/handler/http/`**: Contains HTTP-specific handlers and routing.
          * **`router.go`**: Defines all API routes and maps them to handler functions.
          * **`request.go` & `response.go`**: Defines the Go structs for incoming JSON request bodies and outgoing JSON responses, keeping your handlers clean and type-safe.
          * **`middleware/`**: A dedicated folder for all your HTTP middleware functions.
              * `auth.go`: For JWT validation, token refresh, and user role checks.
              * `ratelimiter.go`: To implement request rate limiting.
              * `logger.go`: For logging details about each incoming request.
          * **Clean Architecture Layer:** **Interface Adapters (Controllers)**.
      * **`internal/service/`**:
          * **Purpose:** Contains clients for interacting with external services that are not direct data persistence layers, such as the AI API for content generation.
          * **Content:** `ai_service.go` will encapsulate calls to your AI integration.
          * **Clean Architecture Layer:** **Frameworks & Drivers**.
  * **`pkg/`**:
      * **Purpose:** This directory is for shared, reusable utility code that is *not specific to your application's business logic*. These packages could theoretically be extracted and used in other Go projects.
      * **Content:**
          * `config/`: A utility to load application configuration.
          * `database/`: Functions to initialize and manage database connections (e.g., MongoDB client setup).
          * `hash/`: A wrapper for password hashing (e.g., bcrypt).
          * `jwt/`: Helper functions to create, sign, and validate JWTs.
          * `logger/`: Setup for a structured logger (e.g., `slog` or a custom wrapper).
          * `validator/`: Common data validation utilities.
  * **`migrations/`**:
      * **Purpose:** Manages database schema changes. For MongoDB (NoSQL), this would typically contain Go scripts to programmatically manage indexes or perform data transformations.
      * **Content:** `000001_create_initial_indexes.go` (example) would contain Go code using the MongoDB driver to ensure necessary indexes (e.g., on `user` emails or `blog` titles) are created for efficient querying.
  * **`.env.example`**: A template for environment variables used by the application.
  * **`.gitignore`**: Specifies intentionally untracked files to ignore.
  * **`go.mod` & `go.sum`**: Go module files for dependency management.
  * **`README.md`**: Project description, setup instructions, and key information.

## 3\. Naming Conventions

Consistent and descriptive naming is paramount for code readability, maintainability, and effective collaboration.

### 3.1. Go General Naming Conventions

  * **Variables (Local & Package-Level):**
      * **Convention:** `camelCase` (e.g., `userName`, `blogPostTitle`).
      * **Constants (Exported):** `UPPER_SNAKE_CASE` (e.g., `MAX_RETRIES`, `DEFAULT_PORT`).
      * **Constants (Unexported):** `camelCase` (e.g., `defaultTimeout`).
      * **Booleans:** Prefix with `is`, `has`, `can` (e.g., `isLoggedIn`, `hasPermission`).
      * **Arrays/Slices:** Use plural forms (e.g., `users`, `posts`).
  * **Functions & Methods:**
      * **Convention:** `PascalCase` for exported functions/methods, `camelCase` for unexported (private) functions/methods.
      * **Examples:** `CreateUser()`, `GetUserByID()`, `calculateHash()`.
      * **Interface Methods:** `PascalCase` (e.g., `Create(ctx context.Context, user *entity.User) error`).
  * **Types (Structs, Interfaces, Custom Types):**
      * **Convention:** `PascalCase` (e.g., `User`, `Blog`, `UserRepository`, `AuthConfig`).
  * **Packages:**
      * **Convention:** All `lowercase`, short, concise, and descriptive. Avoid plural forms.
      * **Examples:** `entity`, `usecase`, `repository`, `handler`, `jwt`, `hash`, `config`.
  * **Files:**
      * **Convention:** `snake_case` or `kebab-case` if not a direct Go package file. For Go source files, `snake_case` is common for multiple logical units (e.g., `user_handler.go`), otherwise the package name is often used (e.g., `jwt.go` for the `jwt` package).
      * **Examples:** `user_usecase.go`, `router.go`, `config.yml`.

### 3.2. MongoDB Specific Naming Conventions

  * **Database Names:**
      * **Convention:** `snake_case` or `kebab-case`. Keep concise.
      * **Example:** `blog_platform_db` or `g6-blog-db`.
  * **Collection Names:**
      * **Convention:** `snake_case` (common) or `camelCase`. **Plural nouns** are strongly recommended for collections of documents.
      * **Examples:** `users`, `posts`, `tokens`.
  * **Document Field Names (within MongoDB):**
      * **Convention:** `camelCase` is idiomatic for JSON/BSON fields.
      * **Go Struct Tags:** Use `bson:"camelCase"` and `json:"camelCase"` tags on your `PascalCase` Go struct fields to ensure proper mapping.
      * **Example:**
          * Go struct field: `CreatedAt`
          * MongoDB document field: `createdAt`
          * Go struct tag: `bson:"createdAt" json:"createdAt"`

### 3.3. Specific Go Entity Naming (`internal/entity/`)

  * **Struct Names:** `PascalCase`, singular noun.
      * **Examples:** `User`, `Blog`, `Token`.
  * **Struct Field Names:** `PascalCase` (exported). Apply `json` and `bson` tags for serialization to `camelCase`.
      * **Example (`user.go`):**
        ```go
        type User struct {
            ID        string    `bson:"_id,omitempty" json:"id"`
            Username  string    `bson:"username" json:"username"`
            Email     string    `bson:"email" json:"email"`
            Password  string    `bson:"password" json:"-"` // Not exposed
            Role      string    `bson:"role" json:"role"`
            CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
            UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
        }
        ```

### 3.4. Specific Go Interface Naming (`internal/usecase/interfaces.go`)

  * **Interface Names:** `PascalCase`, typically `[EntityName]Repository` or `[Service]Service`.
      * **Examples:** `UserRepository`, `BlogRepository`, `TokenRepository`, `AIService`.
  * **Interface Method Names:** `PascalCase`, verb-oriented.
      * **Examples:** `Create(ctx context.Context, user *entity.User) error`, `GetByID(ctx context.Context, id string) (*entity.Blog, error)`, `GenerateContent(ctx context.Context, keywords []string) (string, error)`.

### 3.5. Go Concrete Implementation Naming (e.g., `internal/repository/mongodb/`, `internal/usecase/`)

  * **Structs Implementing Interfaces:** `PascalCase`. If there are multiple implementations for an interface (e.g., `PostgresUserRepository` and `MongoUserRepository`), use a technology prefix. If it's the only implementation, a simpler name like `UserUseCase` is fine.
      * **Examples:** `MongoUserRepository`, `BlogUseCase`, `DefaultAIService`.
  * **Methods (on implementations):** Follow standard Go function naming (PascalCase if exported, camelCase if unexported/helper).
      * **Example (`MongoUserRepository` method):** `Create(ctx context.Context, user *entity.User) error`.

## 4\. GitHub Collaboration Conventions and Workflows

A production-grade GitHub workflow emphasizes clear branching, frequent integration, and thorough code reviews to ensure stability and quality.

### 4.1. Branching Strategy (GitFlow-inspired)

  * **`main` Branch:**
      * **Purpose:** Represents the production-ready, stable codebase.
      * **Rules:** Only updated via merges from `release` or `hotfix` branches. Direct commits are strictly forbidden. Tagged for releases (e.g., `v1.0.0`).
  * **`develop` Branch:**
      * **Purpose:** Integrates new features and bug fixes from feature branches. It should always reflect the latest delivered development changes.
      * **Rules:** Feature branches are merged into `develop`. `release` branches are created from `develop`.
  * **`feature/<feature-name>` Branches:**
      * **Purpose:** For developing new features or significant enhancements.
      * **Rules:** Branch off `develop`. Name descriptively (e.g., `feature/user-authentication`, `feature/blog-search`). Kept short-lived and merged back into `develop` once complete.
  * **`bugfix/<bug-description>` Branches:**
      * **Purpose:** For addressing non-critical bugs discovered during development or QA, or minor fixes not requiring a hotfix.
      * **Rules:** Branch off `develop`. Merged back into `develop`.
  * **`release/<version-number>` Branches:**
      * **Purpose:** For preparing a new production release (e.g., `release/v1.0.0`). Used for final bug fixes, testing, and preparing release notes.
      * **Rules:** Branch off `develop`. Once stable and ready, merged into both `main` and `develop`.
  * **`hotfix/<hotfix-description>` Branches:**
      * **Purpose:** To quickly address critical bugs in the `main` (production) branch.
      * **Rules:** Branch off `main`. Once fixed, merged into both `main` and `develop`.

### 4.2. Pull Request (PR) Workflow

  * **Create PRs Early and Often:** Even for work in progress (use "Draft PRs" or prefix `[WIP]` in title). This encourages early feedback.
  * **Clear and Structured Titles:** Follow Conventional Commits (see below).
      * **Examples:** `feat: Add user profile update endpoint`, `fix: Resolve blog pagination issue`, `refactor(auth): Improve JWT token validation`.
  * **Descriptive PR Descriptions:** Provide comprehensive context:
      * **Problem Solved:** What issue or user story does this PR address?
      * **Solution Implemented:** Briefly explain how the code solves the problem.
      * **Testing Done:** How was this change tested (unit, integration, manual steps)?
      * **Screenshots/GIFs:** Include if applicable for UI changes or complex flows.
      * **Related Issues/Tasks:** Link to Jira, GitHub Issues, or other tracking systems (e.g., `Closes #123`, `Fixes ABC-456`).
      * **Checklist:** A small checklist of done/remaining tasks (`- [x] unit tests added`, `- [ ] documentation updated`).
  * **Keep PRs Small and Focused:** Each PR should ideally address a single feature, bug, or logical task. This makes reviews easier and faster.
  * **Mandatory Code Review:**
      * Require at least one (preferably two) approving reviews from team members.
      * Reviewers focus on code quality, maintainability, adherence to coding standards, potential side effects, performance, and security implications.
      * Use automated tools (linters, static analysis, test results) as pre-merge checks.
  * **Merge Strategy:** Use "Squash and Merge" for feature and bugfix branches into `develop` to maintain a clean, linear history. Use "Merge Pull Request" for `release` and `hotfix` branches into `main` and `develop`.
  * **Branch Cleanup:** Delete merged feature and bugfix branches after they are integrated into `develop`.

## 5\. Commit Conventions

Good commit messages are fundamental for a traceable, understandable, and maintainable codebase. They act as a living documentation of your project's history.

### 5.1. Basic Rules

  * **Atomic Commits:** Each commit should represent a single, logical, and complete change. This means:
      * Fixing one bug per commit.
      * Implementing one small feature or part of a feature per commit.
      * Avoid mixing unrelated changes in a single commit.
      * This makes debugging (`git blame`, `git bisect`) and reverting changes much easier.
  * **Commit Often:** Make small, frequent commits. This reduces the chance of large, complex merge conflicts and allows for more regular integration.
  * **Don't Commit Half-Done Work:** Only commit code when a logical component or a set of related changes is completed and tested. If you need to switch branches or temporarily save changes, use `git stash`.
  * **Test Your Code:** Ensure your code passes all relevant tests (unit, integration) before committing. A broken `main` or `develop` branch is a serious issue.

### 5.2. Conventional Commits (Highly Recommended for Production Grade)

The Conventional Commits specification provides a lightweight convention for standardized commit messages. This enables automated changelog generation, semantic versioning, and clearer communication.

**Format:**

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

  * **`<type>` (Mandatory):** Describes the *nature* of the change. Use lowercase.
      * `feat`: A new feature (correlates with `MINOR` version bump).
          * **Example:** `feat: Add user profile update endpoint`
      * `fix`: A bug fix (correlates with `PATCH` version bump).
          * **Example:** `fix(auth): Resolve token refresh expiration bug`
      * `docs`: Documentation only changes.
          * **Example:** `docs: Update API documentation for blog retrieval`
      * `style`: Changes that do not affect the meaning of the code (whitespace, formatting, semicolons, etc.).
          * **Example:** `style: Format user_handler.go according to gofmt`
      * `refactor`: A code change that neither fixes a bug nor adds a feature (e.g., restructuring code).
          * **Example:** `refactor(blog): Decouple tag processing from creation logic`
      * `perf`: A code change that improves performance.
          * **Example:** `perf: Optimize MongoDB queries for blog search`
      * `test`: Adding missing tests or correcting existing tests.
          * **Example:** `test: Add unit tests for user registration use case`
      * `build`: Changes that affect the build system or external dependencies (e.g., updating Go modules, build scripts).
          * **Example:** `build: Upgrade go.mod dependencies to latest stable versions`
      * `ci`: Changes to CI configuration files and scripts (e.g., GitHub Actions, Jenkins).
          * **Example:** `ci: Add linting step to CI pipeline`
      * `chore`: Other changes that don't modify source code or tests (e.g., updating `.gitignore`).
          * **Example:** `chore: Update .gitignore with IDE files`
      * `revert`: Reverts a previous commit.
          * **Example:** `revert: "feat: Implement experimental caching"`
  * **`[optional scope]`:** Provides context for the change, usually indicating the part of the codebase affected. It is enclosed in parentheses.
      * **Example:** `feat(user-management): Add user promotion functionality`
  * **`<description>` (Mandatory):** A concise, imperative, present-tense summary of the change.
      * **Length:** Keep it short (ideally under 50-72 characters).
      * **Good:** `Add user registration endpoint`
      * **Bad:** `Added user registration endpoint` or `User registration endpoint added`
  * **`[optional body]`:** A more detailed, multi-paragraph explanation of the change.
      * **Purpose:** Explain *what* the change does and *why* it was made, not *how* (the code explains how).
      * **Formatting:** Wrap lines at 72 characters.
      * **Example:**
        ```
        feat(blog): Implement advanced search and filtration

        This commit introduces new endpoints for searching blog posts by title or author
        and filtering by tags, date range, and popularity metrics. It leverages MongoDB's
        text search capabilities and aggregation pipeline for efficient data retrieval.
        ```
  * **`[optional footer(s)]`:** Used for referencing issues or indicating **breaking changes**.
      * **`BREAKING CHANGE:`:** Indicates a breaking API change that requires users of your library to make modifications. This correlates with a `MAJOR` version bump.
          * **Example:** `BREAKING CHANGE: User roles are now stored as integers instead of strings.`
      * **Issue References:** Link to issue tracker items. Use keywords like `Closes #123`, `Fixes #456`, `Resolves JIRA-789`.
      * **Example:**
        ```
        fix: Correct pagination logic for blog retrieval

        Addresses an issue where the pagination was miscalculating
        the total number of pages due to an incorrect count query.

        Closes #55
        ```

### 5.3. Tools for Enforcement

  * **Commitlint:** Integrate `commitlint` (or similar tools) into your Git hooks (e.g., using `husky` for JavaScript projects, or Go-specific pre-commit hooks) to automatically validate commit messages against the Conventional Commits specification before they are committed.

-----

By rigorously following these guidelines, your G6 Blog Starter Project will benefit from a robust, scalable architecture, clear, consistent code, and a highly efficient and transparent development workflow.