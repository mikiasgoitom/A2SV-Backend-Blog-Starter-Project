package middleware

import (
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
)

func RateLimiter(lmt *limiter.Limiter) gin.HandlerFunc {
	return tollbooth_gin.LimitHandler(lmt)
}
