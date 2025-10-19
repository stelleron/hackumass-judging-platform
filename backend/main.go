package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func setupRouter() *gin.Engine {
	// Initialize engine
	r := gin.Default()

	// Set up router
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

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

	// Set up router
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":" + backend_port)
}
