package domain

import (
	"time"
)

type CartItem struct {
	ProductID	string	`json:"productId"`
	Name		string	`json:"name"`
	Price		float64	`json:"price"`
	Quantity	int		`json:"quantity"`
}

type Cart struct {
	UserID		string		`json:"userId"`
	Items		[]CartItem	`json:"items"`
	TotalItem	float64		`json:"totalAmount"`
	UpdatedAt	time.Time	`json:"updatedAt"`
}