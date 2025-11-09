package domain

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"` // Hashed password
	IsAdmin  bool   `gorm:"default:false"`
	Cart     Cart   `gorm:"foreignKey:UserID"`
}