package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// APIHandler holds the business logic and utility services required by the handlers.
type APIHandler struct {
	Service    ECommerceService
	JWTService JWTService
}

// Utility function to respond with JSON
func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Error marshaling JSON response"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// Utility function to respond with JSON error
func RespondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// getUserClaims is a helper to retrieve user data from the request context
func GetUserClaims(r *http.Request) *Claims {
	claims, ok := r.Context().Value(UserContextKey).(*Claims)
	if !ok || claims == nil {
		return nil // Should be caught by middleware, but safety check
	}
	return claims
}

// --- USER AND AUTH HANDLERS ---

// --- ADMIN PRODUCT HANDLERS ---

func (h *APIHandler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Business validation in service, but initial checks here
	if product.Name == "" || product.PriceCents <= 0 || product.Inventory < 0 {
		respondError(w, http.StatusBadRequest, "Product name, price, and inventory must be valid")
		return
	}

	if err := h.Service.CreateProduct(&product); err != nil {
		respondError(w, http.StatusInternalServerError, "Could not create product")
		return
	}
	respondJSON(w, http.StatusCreated, product)
}

// --- PUBLIC PRODUCT HANDLERS ---

func (h *APIHandler) GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	products, err := h.Service.GetProducts(query)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to retrieve products")
		return
	}

	// Format price for display (Price is stored in cents, display in dollars)
	displayProducts := make([]map[string]interface{}, len(products))
	for i, p := range products {
		displayProducts[i] = map[string]interface{}{
			"id":          p.ID,
			"name":        p.Name,
			"description": p.Description,
			"price":       fmt.Sprintf("%.2f", float64(p.PriceCents)/100.0), // Convert cents to dollars
			"inventory":   p.Inventory,
		}
	}
	respondJSON(w, http.StatusOK, displayProducts)
}

// --- CART HANDLERS (Authenticated) ---

func (h *APIHandler) AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	claims := getUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "User context missing")
		return
	}

	var req struct {
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Quantity <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid product ID or quantity")
		return
	}

	cart, err := h.Service.AddToCart(claims.UserID, req.ProductID, req.Quantity)
	if err != nil {
		if errors.Is(err, ErrInsufficientInv) || errors.Is(err, ErrNotFound) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to add item to cart")
		return
	}

	respondJSON(w, http.StatusOK, cart)
}

func (h *APIHandler) ViewCartHandler(w http.ResponseWriter, r *http.Request) {
	claims := getUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "User context missing")
		return
	}

	cart, err := h.Service.ViewCart(claims.UserID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to retrieve cart")
		return
	}

	// Calculate and format total for frontend display
	totalCents := int64(0)
	for _, item := range cart.Items {
		totalCents += item.PriceCents * int64(item.Quantity)
	}

	response := map[string]interface{}{
		"cart":        cart,
		"total_usd":   fmt.Sprintf("%.2f", float64(totalCents)/100.0),
		"total_cents": totalCents,
	}
	respondJSON(w, http.StatusOK, response)
}

func (h *APIHandler) CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	claims := getUserClaims(r)
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "User context missing")
		return
	}

	result, err := h.Service.Checkout(claims.UserID)
	if err != nil {
		if errors.Is(err, ErrCartEmpty) || errors.Is(err, ErrInsufficientInv) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		// Generic internal server error for payment gateway issues
		respondError(w, http.StatusInternalServerError, "Checkout failed due to internal error.")
		return
	}

	respondJSON(w, http.StatusOK, result)
}
