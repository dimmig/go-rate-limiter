package middleware

import (
	"fmt"
	"github.com/dimmig/go-rate-limiter/repository"
	"github.com/dimmig/go-rate-limiter/service"
	"github.com/gin-gonic/gin"
)

type RateLimiterConfig struct {
	Limiter *service.LimiterService
	KeyFunc func(ctx *gin.Context) string
	Limit   int
	Expiry  int
}

func NewRateLimiterConfig(keyFunc func(ctx *gin.Context) string, limit, expiry int) RateLimiterConfig {
	repo := repository.NewRedisRepo()
	limiter := service.NewLimiterService(repo)
	return RateLimiterConfig{
		Limiter: limiter,
		KeyFunc: keyFunc,
		Limit:   limit,
		Expiry:  expiry,
	}
}

func NewRateLimiter(cfg RateLimiterConfig) gin.HandlerFunc {
	return func(g *gin.Context) {
		key := cfg.KeyFunc(g)

		allowed, remaining, err := cfg.Limiter.Allow(key, cfg.Limit, cfg.Expiry)
		if err != nil {
			g.JSON(500, gin.H{"error": "rate limiter error"})
			g.Abort()
			return
		}
		if !allowed {
			g.JSON(429, gin.H{"error": "rate limit exceeded"})
			g.Abort()
			return
		}
		g.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

		g.Next()
	}
}
