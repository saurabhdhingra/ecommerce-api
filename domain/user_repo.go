package domain

type UserRepository interface {
	Create(user *User) error
	FindByUsername(username string) (*User, error)
	FindByID(id uint) (*User, error)
}