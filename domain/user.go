package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleCustomer	UserRole = "customer"
	RoleAdmin		UserRole = "admin"
)

type User struct {
	ID			string		`json:"id"`
	Username	string		`json:"username"`
	Email		string		`json:"email"`
	Password	string		`json:"-"`
	UserRole	UserRole	`json:"role"`
	CreatedAt	time.Time	`json:"createdAt"`
}

type UserClaims struct {
	ID		string		`json:"id"`
	Email	string		`json:"email"`
	Role 	UserRole	`json:"role"`
}

func NewUserID() string {
	return uuid.New().String()
}