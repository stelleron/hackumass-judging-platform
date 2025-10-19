package middleware

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Hash the password before storing it
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

// Verify the given password with the stored one
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	if err != nil {
		return false, "invalid credentials"
	}
	return true, ""
}
