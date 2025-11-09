package domain

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string  `gorm:"not null"`
	Description string
	PriceCents  int64   `gorm:"not null"` // Price stored in smallest currency unit (cents)
	Inventory   int     `gorm:"default:0"`
}