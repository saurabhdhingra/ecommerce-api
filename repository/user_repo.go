package repository

import (
	"gorm.io/gorm"
	
	"ecommerce-api/domain"
)


func (r *PostgresRepository) CreateUser(user *domain.User) error {
	return r.DB.Create(user).Error
}

func (r *PostgresRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, domain.ErrNotFound
	}
	return &user, err
}

func (r *PostgresRepository) FindUserByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.DB.First(&user, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, domain.ErrNotFound
	}
	return &user, err
}