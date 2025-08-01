package contract

import "github.com/golang-jwt/jwt/v5"

type IJWTManager interface {
	GenerateAccessToken(userID, userRole string) (string, error)
	GenerateRefreshToken(tokenID, userID string) (string, error)
	VerifyToken(token string) (*CustomClaims, error)
	VerifyRefreshToken(token string) (*RefreshClaims, error)
}

type CustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	jwt.RegisteredClaims
}
