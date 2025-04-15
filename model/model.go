package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique"` 
	Roll     string `json:"Roll"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
type USer_Delete struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
type Class struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      uint   `json:"user_id"` // FK to User.ID
	User        User   `gorm:"foreignKey:UserID"`
}
type Class_delete struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
