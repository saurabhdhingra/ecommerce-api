package repository

import (
	"gorm.io/gorm"

	"ecommerce-api/domain"
)

type ProductRepo struct {
	*PostgresRepository
}

func (r *ProductRepo) Create(product *domain.Product) error {
	return r.DB.Create(product).Error
}

func (r *ProductRepo) FindAll(query string) ([]domain.Product, error) {
	var products []domain.Product
	db := r.DB
	if query != "" {
		db = db.Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%")
	}
	err := db.Find(&products).Error
	return products, err
}

func (r *ProductRepo) FindByID(id uint) (*domain.Product, error) {
	var product domain.Product
	err := r.DB.First(&product, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, domain.ErrNotFound
	}
	return &product, err
}

func (r *ProductRepo) Update(product *domain.Product) error {
	return r.DB.Save(product).Error
}

func (r *ProductRepo) Delete(id uint) error {
	err := r.DB.Delete(&domain.Product{}, id).Error
	if err != nil {
		return err
	}
	return nil
}
