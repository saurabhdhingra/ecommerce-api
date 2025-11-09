

func (s *serviceImpl) CreateProduct(product *Product) error {
	return s.productRepo.Create(product)
}

func (s *serviceImpl) GetProducts(query string) ([]Product, error) {
	return s.productRepo.FindAll(query)
}