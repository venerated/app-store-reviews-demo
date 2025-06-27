package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// CORS settings
	router.Use(cors.New(cors.Config{
		// Allow React app access
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/reviews", handleReviewsRequest)
	router.Run("localhost:8080")
}
