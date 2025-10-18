// examples/04_different_routes.go
// Apply different rate limits to different route groups
// Usage: go run 04_different_routes.go
// Test public: curl http://localhost:8080/public/health
// Test API: curl http://localhost:8080/api/data

package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"go-rate-limiter/middleware"
)

func main() {
	r := gin.Default()

	publicLimiter := middleware.NewRateLimiterConfig(
		func(c *gin.Context) string {
			return c.ClientIP()
		},
		5,  // Only 5 requests per 10 seconds
		10, // for public endpoints
	)

	authLimiter := middleware.NewRateLimiterConfig(
		func(c *gin.Context) string {
			userID := c.GetHeader("X-User-ID")
			if userID == "" {
				return c.ClientIP()
			}
			return userID
		},
		100, // 100 requests per 60 seconds
		60,  // for authenticated users
	)

	// PUBLIC ROUTES - Strict limit
	publicGroup := r.Group("/public")
	publicGroup.Use(middleware.NewRateLimiter(publicLimiter))
	{
		publicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		publicGroup.POST("/login", func(c *gin.Context) {
			c.JSON(200, gin.H{"token": "abc123xyz"})
		})

		publicGroup.GET("/info", func(c *gin.Context) {
			c.JSON(200, gin.H{"info": "public information"})
		})
	}

	// AUTHENTICATED ROUTES - Relaxed limit
	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.NewRateLimiter(authLimiter))
	{
		apiGroup.GET("/data", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"data":    "sensitive data",
				"user":    c.GetHeader("X-User-ID"),
				"records": 1000,
			})
		})

		apiGroup.POST("/data", func(c *gin.Context) {
			c.JSON(201, gin.H{
				"message": "data created",
				"id":      12345,
			})
		})

		apiGroup.GET("/stats", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"requests": 5000,
				"users":    150,
			})
		})
	}

	log.Println("Server running on :8080")

	r.Run(":8080")
}
