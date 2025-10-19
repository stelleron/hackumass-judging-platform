package router

import (
	"context"
	"fmt"
	"net/http"
	"time"

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
	// Check if request valid
	var creds models.Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Hash password and add to MongoDB
	hash := middleware.HashPassword(creds.Password)
	userCollection := Client.Database("hackumass-xiii").Collection("credentials")
	_, err := userCollection.InsertOne(context.TODO(), bson.M{
		"username": creds.Username,
		"password": string(hash),
	})

	// If user not created, return error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}

	// Return json
	c.JSON(http.StatusOK, gin.H{"message": "user created"})
}

func Login(c *gin.Context) {
	// Check if request valid
	var creds models.Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	fmt.Println("Login attempt:", creds.Username)

	// Get same credentials from MongoDB
	var result bson.M
	userCollection := Client.Database("hackumass-xiii").Collection("credentials")
	err := userCollection.FindOne(context.TODO(), bson.M{"username": creds.Username}).Decode(&result)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	fmt.Println("DB result:", result)

	// Check passwords
	check, msg := middleware.VerifyPassword(result["password"].(string), creds.Password)
	if !check {
		c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
		return
	}

	fmt.Println("Password hash in DB:", result["password"])

	// Set JWT for 1 day
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)

	// Create cookie and return OK
	c.SetCookie("token", tokenString, 3600*24, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}

func Verify(c *gin.Context) {
	tokenStr, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !tkn.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": claims.Username})
}

// DEBUGGING ONLY
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
