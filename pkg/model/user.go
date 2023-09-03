package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName    string `gorm:"size:25;default:null;not null"`
	LastName     string `gorm:"size:25;default:null;not null"`
	Email        string `gorm:"unique"`
	PasswordHash string `gorm:"default:null;not null"`
}

type RegisterDTO struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=25"`
	LastName  string `json:"last_name" validate:"required,min=2,max=25"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}
