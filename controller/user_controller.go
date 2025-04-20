package controller

import (
	"GOLANG/model"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}
func GenerateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24h
	}
	jwtKey := os.Getenv("jwt_key")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
func (uc *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []model.User
	if err := uc.DB.Find(&users).Error; err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser model.User

	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var existingUser model.User
	if err := uc.DB.Where("username = ?", newUser.Username).First(&existingUser).Error; err == nil {
		// If a user is found with the same username, return a 409 conflict
		http.Error(w, "Username already taken", http.StatusConflict)
		return
	} else if err != gorm.ErrRecordNotFound {
		http.Error(w, "Failed to check username", http.StatusInternalServerError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	newUser.Password = string(hashedPassword)

	// Create the new user in the database
	if err := uc.DB.Create(&newUser).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Return success message
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User creation complete"})
}

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var input model.USer_Delete
	// Decode JSON body into `input`
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Check if username or email was provided
	if input.Username == "" && input.Email == "" {
		http.Error(w, "Username or email required", http.StatusBadRequest)
		return
	}

	var user model.User
	var result *gorm.DB

	// Find user by username or email
	if input.Username != "" {
		result = uc.DB.Where("username = ?", input.Username).First(&user)
	} else {
		result = uc.DB.Where("email = ?", input.Email).First(&user)
	}

	if result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Delete user
	if err := uc.DB.Delete(&user).Error; err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}
func (uc *UserController) DeleteAllUsers(w http.ResponseWriter, r *http.Request) {
	// Delete all rows in users table
	if err := uc.DB.Where("1 = 1").Delete(&model.User{}).Error; err != nil {
		http.Error(w, "Failed to delete all users", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "All users have been deleted successfully",
	})
}
func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var creds model.User_login
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var user model.User
	if err := uc.DB.Where("username = ?", creds.Username).First(&user).Error; err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := GenerateJWT(user.Username)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

		logFile, err := os.OpenFile("logins.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			logger := log.New(logFile, "", log.LstdFlags)
			logger.Printf("User '%s' logged in at %s", creds.Username, time.Now().Format(time.RFC3339))
			defer logFile.Close()
		}
		
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}
func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	// Overwrite the token cookie with expired value
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		HttpOnly: true,
		Secure:   false, // Set to true in production
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0), // Expired
	})
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}