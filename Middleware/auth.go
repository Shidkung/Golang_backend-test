package middleware

import (
	"fmt"
	"net/http"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"os"
	"log"
)
var JwtKey []byte

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the JWT secret key from the environment variables
	JwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the "Authorization" cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			// Handle case where cookie is missing
			http.Error(w, "Unauthorized - missing cookie", http.StatusUnauthorized)
			return
		}

		// The cookie contains the JWT token
		tokenStr := cookie.Value

		// Log the token being parsed for debugging
		fmt.Println("Received Token:", tokenStr)

		// Parse and validate the token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Check token signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return JwtKey, nil
		})

		if err != nil || !token.Valid {
			// Log the error for debugging
			fmt.Println("Error Parsing Token:", err)
			http.Error(w, "Unauthorized - invalid token", http.StatusUnauthorized)
			return
		}

		// Optionally, you can store claims (like username) into the context for use later
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized - invalid claims", http.StatusUnauthorized)
			return
		}

		// Save username or other claims into context
		r = r.WithContext(context.WithValue(r.Context(), "username", claims["username"]))

		// Call next handler
		next(w, r)
	}
}
