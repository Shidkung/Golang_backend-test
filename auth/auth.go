package auth

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var jwtKey []byte

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Fetch JWT_SECRET_KEY from environment and assign it to jwtKey
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		log.Fatal("JWT_SECRET_KEY not set in the environment")
	} else {
		jwtKey = []byte(secret) // Assign secret to jwtKey variable
		log.Println("JWT_SECRET_KEY is set")
	}
}

// GenerateJWT creates a new JWT for the user
func GenerateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey) // Signs and returns the JWT token
}
