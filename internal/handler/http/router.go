package http

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, authHandler AuthHandler) {
	r.GET("/auth/google/login", authHandler.HandleGoogleLogin)
	r.GET("/auth/google/callback", authHandler.HandleGoogleCallback)
}
