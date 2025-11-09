package domain


type CartRepository interface {
	FindByUserID(userID uint) (*Cart, error)
	Save(cart *Cart) error
	Clear(userID uint) error
	UpdateInventory(productID uint, quantityChange int) error // For managing inventory during checkout
}