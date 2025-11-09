package domain

type ProductRepository interface {
	Create(product *Product) error
	FindAll(query string) ([]Product, error)
	FindByID(id uint) (*Product, error)
	Update(product *Product) error
	Delete(id uint) error
}