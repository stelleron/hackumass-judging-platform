package main

import (
	"log"
	"os"

	"hackumass-xiii.com/judging-platform/router"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func setupRouter(mongodb_url string) *gin.Engine {
	// Initialize engine
	r := gin.Default()

	// Allow CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // frontend origin
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize database
	router.InitMongoDB(mongodb_url)

	// Set up router
	api := r.Group("/api")
	{
		api.POST("/login", router.Login)
		api.POST("/signup", router.Signup)
		api.GET("/verify", router.Verify)
	}

	// Return
	return r
}

func main() {
	// Get .env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	backend_port := os.Getenv("BACKEND_PORT")
	mongodb_url := os.Getenv("MONGODB_URL")

	// Set up router
	r := setupRouter(mongodb_url)

	// Listen and Server in 0.0.0.0:8080
	r.Run(":" + backend_port)
}
