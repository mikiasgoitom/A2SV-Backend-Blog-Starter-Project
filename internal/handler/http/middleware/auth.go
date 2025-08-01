package middleware

import (
	"blog/internal/infrastructure/jwt"
	"blog/internal/usecase"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleWare(jwtMgr jwt.JWTManager, userUseCase usecase.UserUseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower((parts[0])) != "bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}
		tokenString := parts[1]

		claims, err := jwtMgr.VerifyToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		ctx.Set("userID", claims.Subject)
		ctx.Set("userRole", claims.Role)

		ctx.Next()
	}
}
