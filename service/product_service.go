package service

import (
	"ecommerce-api/domain"
)

func (s *ServiceImpl) CreateProduct(product *domain.Product) error {
	return s.productRepo.Create(product)
}

func (s *ServiceImpl) GetProducts(query string) ([]domain.Product, error) {
	return s.productRepo.FindAll(query)
}
