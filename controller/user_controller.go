package controller

import (
	"GOLANG/model"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
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
