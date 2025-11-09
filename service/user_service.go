package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"

	"ecommerce-api/domain"
)

// ECommerceService defines all usecase operations for the application.
type ECommerceService interface {
	// Auth & User
	Signup(username, password string) (*domain.User, error)
	Login(username, password string) (*domain.User, error)

	// Products
	CreateProduct(product *domain.Product) error
	GetProducts(query string) ([]domain.Product, error)

	// Cart & Checkout
	AddToCart(userID uint, productID uint, quantity int) (*domain.Cart, error)
	RemoveFromCart(userID uint, productID uint, quantity int) (*domain.Cart, error)
	ViewCart(userID uint) (*domain.Cart, error)
	Checkout(userID uint) (map[string]interface{}, error)
}

type ServiceImpl struct {
	ECommerceService
	userRepo    domain.UserRepository
	productRepo domain.ProductRepository
	cartRepo    domain.CartRepository
	stripeSvc   StripeService
}

func NewECommerceService(u domain.UserRepository, p domain.ProductRepository, c domain.CartRepository, s StripeService) ECommerceService {
	return &ServiceImpl{userRepo: u, productRepo: p, cartRepo: c, stripeSvc: s}
}

// hashPassword is a simple utility (use bcrypt in production!)
func hashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

// --- Auth & User ---

func (s *ServiceImpl) Signup(username, password string) (*domain.User, error) {
	if username == "" || password == "" {
		return nil, errors.New("username and password cannot be empty")
	}

	hashedPassword := hashPassword(password)
	user := &domain.User{
		Username: username,
		Password: hashedPassword,
		IsAdmin:  false,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err // Handle unique constraint errors in handler/repository
	}
	return user, nil
}

func (s *ServiceImpl) Login(username, password string) (*domain.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, domain.ErrInvalidCredentials // Hide specific error for security
	}

	hashedInput := hashPassword(password)
	// In production, use bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if user.Password != hashedInput {
		return nil, domain.ErrInvalidCredentials
	}
	return user, nil
}

func (s *ServiceImpl) AddToCart(userID uint, productID uint, quantity int) (*domain.Cart, error) {
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	if product.Inventory < quantity {
		return nil, domain.ErrInsufficientInv
	}

	cart, _ := s.cartRepo.FindByUserID(userID)

	found := false
	for i, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items[i].Quantity += quantity
			found = true
			break
		}
	}

	if !found {
		cart.Items = append(cart.Items, domain.CartItem{
			ProductID:  productID,
			Quantity:   quantity,
			Name:       product.Name,
			PriceCents: product.PriceCents,
		})
	}

	if err := s.cartRepo.Save(cart); err != nil {
		return nil, err
	}

	return cart, nil
}

func (s *ServiceImpl) RemoveFromCart(userID uint, productID uint, quantity int) (*domain.Cart, error) {
	cart, err := s.cartRepo.FindByUserID(userID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	newItems := []domain.CartItem{}
	removed := false

	for _, item := range cart.Items {
		if item.ProductID == productID {
			if item.Quantity > quantity {
				item.Quantity -= quantity
				newItems = append(newItems, item)
			}
			removed = true
		} else {
			newItems = append(newItems, item)
		}
	}

	if !removed {
		return nil, errors.New("product not found in cart")
	}

	cart.Items = newItems
	if err := s.cartRepo.Save(cart); err != nil {
		return nil, err
	}

	return cart, nil
}

func (s *ServiceImpl) ViewCart(userID uint) (*domain.Cart, error) {
	return s.cartRepo.FindByUserID(userID)
}

func (s *ServiceImpl) Checkout(userID uint) (map[string]interface{}, error) {
	cart, err := s.cartRepo.FindByUserID(userID)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	if len(cart.Items) == 0 {
		return nil, domain.ErrCartEmpty
	}

	var totalAmount int64

	for _, item := range cart.Items {
		if err := s.cartRepo.UpdateInventory(item.ProductID, -item.Quantity); err != nil {
			log.Printf("Inventory failure for product %d: %v", item.ProductID, err)
			return nil, domain.ErrInsufficientInv
		}
		totalAmount += item.PriceCents * int64(item.Quantity)
	}

	pi, err := s.stripeSvc.CreatePaymentIntent(totalAmount, "usd", "E-commerce order from user")
	if err != nil {
		return nil, errors.New("payment gateway failed to create intent")
	}

	if err := s.cartRepo.Clear(userID); err != nil {
		log.Printf("Warning: Payment succeeded but failed to clear cart for user %d: %v", userID, err)
	}

	return map[string]interface{}{
		"message":           "Checkout successful. Payment initiated.",
		"total_paid_cents":  totalAmount,
		"payment_intent_id": pi.ID,
		"client_secret":     pi.ClientSecret,
	}, nil
}
