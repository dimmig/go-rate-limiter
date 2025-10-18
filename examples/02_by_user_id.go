// examples/02_by_user_id.go
// Rate limiting by user ID (extracted from JWT)
// Usage: go run 02_by_user_id.go
// Test: curl -H "X-User-ID: user123" http://localhost:8080/api/posts

package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"go-rate-limiter/middleware"
)

func main() {
	r := gin.Default()

	// Middleware to extract user ID from header
	// In production, parse JWT token instead
	r.Use(func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			userID = "anonymous"
		}
		c.Set("user_id", userID)
		c.Next()
	})

	rateLimitMW := middleware.NewRateLimiterConfig(
		func(c *gin.Context) string {
			userID, _ := c.Get("user_id")
			return userID.(string)
		},
		10,
		60,
	)

	r.Use(middleware.NewRateLimiter(rateLimitMW))

	r.GET("/api/posts", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		c.JSON(200, gin.H{
			"user":  userID,
			"posts": []string{"post1", "post2", "post3"},
		})
	})

	r.GET("/api/profile", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		c.JSON(200, gin.H{
			"user":  userID,
			"name":  "John Doe",
			"email": "john@example.com",
		})
	})

	log.Println("Server running on :8080")

	r.Run(":8080")
}
