package controller

import (
	"GOLANG/model"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

type Classes_controller struct {
	DB *gorm.DB
}

func (uc *Classes_controller) GetClass(w http.ResponseWriter, r *http.Request) {
	var users []model.Class
	if err := uc.DB.Find(&users).Error; err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (uc *Classes_controller) CreateClass(w http.ResponseWriter, r *http.Request) {
	var newUser model.Class
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if err := uc.DB.Create(&newUser).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User creation complete"})
}

func (uc *Classes_controller) DeleteClass(w http.ResponseWriter, r *http.Request) {
	var input model.Class_delete
	// Decode JSON body into `input`
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Check if username or email was provided
	if input.Title == "" && input.Description == "" {
		http.Error(w, "Username or email required", http.StatusBadRequest)
		return
	}

	var Class model.Class
	var result *gorm.DB

	// Find user by username or email
	if input.Title != "" {
		result = uc.DB.Where("username = ?", input.Title).First(&Class)
	} else {
		result = uc.DB.Where("email = ?", input.Description).First(&Class)
	}

	if result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Delete user
	if err := uc.DB.Delete(&Class).Error; err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}
func (uc *Classes_controller) DeleteAllClasses(w http.ResponseWriter, r *http.Request) {
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
