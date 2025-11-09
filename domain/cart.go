package domain

import (
	"gorm.io/gorm"
)

type CartItem struct {
	gorm.Model
	CartID    uint
	ProductID uint
	Quantity  int `gorm:"default:1"`
	Name string 
	PriceCents int64
}

type Cart struct {
	gorm.Model
	UserID  uint      `gorm:"unique;not null"`
	Items []CartItem
}
