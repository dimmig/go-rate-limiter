// examples/01_basic.go
// Basic rate limiting by IP address
// Usage: go run 01_basic.go
// Test: curl http://localhost:8080/api/data

package main

import (
	"github.com/gin-gonic/gin"
	"go-rate-limiter/middleware"
	"log"
)

func main() {
	r := gin.Default()

	// Create rate limiter
	// Limit: 5 requests per 10 seconds per IP
	rateLimitMW := middleware.NewRateLimiterConfig(
		func(c *gin.Context) string {
			return c.ClientIP()
		},
		5,
		60,
	)

	r.Use(middleware.NewRateLimiter(rateLimitMW))

	r.GET("/api/data", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "success",
			"data":    "hello world",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Println("Server running on :8080")

	r.Run(":8080")
}
