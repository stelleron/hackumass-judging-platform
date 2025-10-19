package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"hackumass-xiii.com/judging-platform/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func setupRouter(mongodb_url string) *gin.Engine {
	// Initialize engine
	r := gin.Default()

	// Initialize MongoDB
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongodb_url).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Send a ping to confirm a successful connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	// Set up router
	api := r.Group("/api")
	{
		/*
			api.POST("/login", handlers.Login)
			api.POST("/signup", handlers.Signup)
		*/
		api.GET("/ping", middleware.Ping)
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
