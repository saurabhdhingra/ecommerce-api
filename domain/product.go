package domain

import (
	"time"
)

type Product struct {
	ID			string		`json:"id"`
	Name		string		`json:"name"`
	Description	string		`json:"description"`
	Price		float64		`json:"price"`
	Inventory	int			`json:"inventory"`
	IsActive	bool		`json:"isActive"`
	CreatedAt	time.Time	`json:"createdAt"`
}
