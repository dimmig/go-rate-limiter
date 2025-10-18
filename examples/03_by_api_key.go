// examples/03_by_api_key.go
// Rate limiting by API key from header
// Usage: go run 03_by_api_key.go
// Test: curl -H "X-API-Key: key123" http://localhost:8080/api/webhook

package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"go-rate-limiter/middleware"
)

func main() {
	r := gin.Default()

	rateLimitMW := middleware.NewRateLimiterConfig(
		func(c *gin.Context) string {
			apiKey := c.GetHeader("X-API-Key")
			if apiKey == "" {
				// No API key provided, rate limit by IP
				return c.ClientIP()
			}
			return "apikey_" + apiKey
		},
		10,
		60,
	)

	r.Use(middleware.NewRateLimiter(rateLimitMW))

	r.POST("/api/webhook", func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(401, gin.H{"error": "X-API-Key header required"})
			return
		}

		c.JSON(200, gin.H{
			"message": "webhook processed",
			"api_key": apiKey,
		})
	})

	r.GET("/api/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   "2024-01-01T00:00:00Z",
		})
	})

	log.Println("Server running on :8080")

	r.Run(":8080")
}
