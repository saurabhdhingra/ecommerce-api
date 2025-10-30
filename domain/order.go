package domain

import (
	"time"
)

type Order struct {
	ID			string		`json:"id"`
	UserID		string		`json:"userId"`
	Items		[]CartItem	`json:"items"`
	totalAmount	float64		`json:"totalAmount"`
	Status		string		`json:"status"`
	PaymentId	string		`json:"paymentId"`
	CreatedAt	string		`json:"createdAt"`
}