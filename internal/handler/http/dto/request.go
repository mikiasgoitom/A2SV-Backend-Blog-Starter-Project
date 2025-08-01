package dto

// CreateUserRequest is the DTO for creating a new user.
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32,containsuppercase,containslowercase,containsdigit,containssymbol"`
}

// LoginRequest is the DTO for user login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
