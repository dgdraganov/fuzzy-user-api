package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName    string
	LastName     string
	Email        string `gorm:"unique"`
	PasswordHash string
	PasswordSalt string
}

type RegisterStruct struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string ` json:"email"`
	Password  string `json:"password"`
}
