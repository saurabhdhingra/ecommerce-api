package repository

import (
	"gorm.io/gorm"

	"ecommerce-api/domain"
)

type UserRepo struct {
	*PostgresRepository
}

func (r *UserRepo) Create(user *domain.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepo) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, domain.ErrNotFound
	}
	return &user, err
}

func (r *UserRepo) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.DB.First(&user, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, domain.ErrNotFound
	}
	return &user, err
}
