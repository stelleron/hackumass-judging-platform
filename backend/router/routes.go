package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"hackumass-xiii.com/judging-platform/middleware"
	"hackumass-xiii.com/judging-platform/models"
)

var Client *mongo.Client // set this in your db init

func InitMongoDB(mongodb_url string) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongodb_url).SetServerAPIOptions(serverAPI)

	var err error
	Client, err = mongo.Connect(context.TODO(), opts) // use =, not :=
	if err != nil {
		panic(err)
	}

	// Ping the DB to verify connection
	if err := Client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}

// User authentication routes
// ====
var jwtKey = []byte("supersecretkey")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Signup(c *gin.Context) {
	var creds models.Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	hash := middleware.HashPassword(creds.Password)
	userCollection := Client.Database("hackumass-xiii").Collection("credentials")
	_, err := userCollection.InsertOne(context.TODO(), bson.M{
		"username": creds.Username,
		"password": string(hash),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user created"})
}

func Login(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func Verify(c *gin.Context) {

}

// DEBUGGING ONLY
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
