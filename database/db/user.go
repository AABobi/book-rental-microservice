package db

import (
	"gorm.io/gorm"
)

type User struct {
	UserID   uint   `json:"user_id" binding:"required" gorm:"primary_key"`
	Email    string `json:"email" binding:"required" gorm:"unique;not null"`
	Password string `json:"password" binding:"required" gorm:"not null"`
}

func (u *User) GetUser(db *gorm.DB) {
	db.Where("email = ?", u.Email).First(&u)
}

func (u *User) AddUser(db *gorm.DB) {
	db.Create(&u)
}

func GetAllUsers(db *gorm.DB, users *[]User) {
	db.Find(&users)
}
